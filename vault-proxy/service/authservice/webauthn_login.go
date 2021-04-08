package authservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) StartLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("auth.username", body.Name))

	resp, err := s.Vault.Logical().Write(fmt.Sprintf("auth/%s/assertion", s.AuthPath), map[string]interface{}{
		"name": body.Name,
	})
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
	sess.Values["login_data"] = resp.Data["session_data"]
	if err := sess.Save(r, w); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", resp.Data["assertion"].(string))
}

func (s *Server) FinishLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	var body struct {
		Name      string          `json:"name"`
		Assertion json.RawMessage `json:"assertion"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("auth.username", body.Name))

	sess, err := s.Store.Get(r, "vault-proxy")
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusBadRequest)
		return
	}

	sessionDataRaw, ok := sess.Values["login_data"]
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

	secret, err := s.Vault.Logical().Write(fmt.Sprintf("auth/%s/login", s.AuthPath), map[string]interface{}{
		"name":               body.Name,
		"session_data":       sessionData,
		"assertion_response": string(body.Assertion),
	})
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusForbidden)
		return
	}

	vaultToken, err := secret.TokenID()
	if err != nil {
		http.Error(w, "no vault token", http.StatusForbidden)
		return
	}
	tokenTTL, err := secret.TokenTTL()
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusForbidden)
		return
	}
	span.SetAttributes(attribute.Float64("auth.token_ttl", tokenTTL.Seconds()))
	tokenAccessor, err := secret.TokenAccessor()
	if err != nil {
		span.RecordError(err)
	} else {
		span.SetAttributes(attribute.String("auth.token_accessor", tokenAccessor))
	}

	sess, err = s.Store.Get(r, "vault-token")
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}
	sess.Options.MaxAge = int(tokenTTL.Seconds())
	sess.Options.Domain = s.CookieDomain
	sess.Options.Path = "/"
	sess.Values["token"] = vaultToken
	sess.Values["creation_ttl"] = int(tokenTTL.Seconds())

	if err := sess.Save(r, w); err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
