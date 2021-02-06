package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) StartRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	span.SetAttributes(label.String("auth.username", body.Name))

	token := r.Header.Get("X-Vault-Token")
	if token == "" {
		http.Error(w, "no X-Vault-Token header found", http.StatusBadRequest)
		return
	}

	vault, err := s.Vault.Clone()
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}
	vault.SetToken(token)

	resp, err := vault.Logical().Read(fmt.Sprintf("auth/%s/users/%s/credentials/request", s.AuthPath, body.Name))
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}

	sess, err := s.Store.Get(r, "vault-proxy")
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}

	// Session should only last 5 more minutes
	sess.Options.MaxAge = 300
	sess.Values["registration_data"] = resp.Data["session_data"]
	if err := sess.Save(r, w); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", resp.Data["creation_response"].(string))
}

func (s *Server) FinishRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	var body struct {
		Name        string          `json:"name"`
		Attestation json.RawMessage `json:"attestation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	span.SetAttributes(label.String("auth.username", body.Name))

	token := r.Header.Get("X-Vault-Token")
	if token == "" {
		http.Error(w, "no X-Vault-Token header found", http.StatusBadRequest)
		return
	}

	sess, err := s.Store.Get(r, "vault-proxy")
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	sessionDataRaw, ok := sess.Values["registration_data"]
	if !ok {
		http.Error(w, "no previous registration session saved", http.StatusForbidden)
		return
	}
	sessionData := sessionDataRaw.(string)

	// delete the session
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	vault, err := s.Vault.Clone()
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}
	vault.SetToken(token)

	_, err = vault.Logical().Write(fmt.Sprintf("auth/%s/users/%s/credentials/create", s.AuthPath, body.Name), map[string]interface{}{
		"session_data":         sessionData,
		"attestation_response": string(body.Attestation),
	})
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
