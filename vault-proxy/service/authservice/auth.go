package authservice

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	vaultTokenHeader := r.Header.Get("X-Vault-Token")
	if vaultTokenHeader != "" {
		span.SetAttributes(attribute.String("auth.token_source", "header"))
		s.handleVaultToken(ctx, r, w, vaultTokenHeader, nil)
		return
	}

	sess, err := s.Store.Get(r, "vault-token")
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}
	if vaultTokenRaw, ok := sess.Values["token"]; ok {
		span.SetAttributes(attribute.String("auth.token_source", "session"))
		s.handleVaultToken(ctx, r, w, vaultTokenRaw.(string), sess)
		return
	}

	http.Error(w, "no credentials present", http.StatusUnauthorized)
}

func (s *Server) handleVaultToken(ctx context.Context, r *http.Request, w http.ResponseWriter, token string, sess *sessions.Session) {
	span := trace.SpanFromContext(ctx)

	vault, err := s.Vault.Clone()
	if err != nil {
		spanerr.RecordError(ctx, err)
		http.Error(w, "failed to clone vault client", http.StatusInternalServerError)
		return
	}

	vault.SetToken(token)
	secret, err := vault.Auth().Token().LookupSelf()
	if err != nil {
		spanerr.RecordError(ctx, err)
		http.Error(w, "invalid vault token", http.StatusForbidden)
		return
	}

	tokenAccessor, err := secret.TokenAccessor()
	if err != nil {
		span.RecordError(err)
	} else {
		span.SetAttributes(attribute.String("auth.token_accessor", tokenAccessor))
	}

	tokenTTL, err := secret.TokenTTL()
	if err != nil {
		http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
		return
	}
	span.SetAttributes(attribute.Stringer("auth.token_ttl", tokenTTL))

	// TODO maybe parameterize this
	if sess != nil && tokenTTL < (24*time.Hour) {
		creationTTL, ok := sess.Values["creation_ttl"]
		if !ok {
			creationTTL = 48 * 3600
		}

		secret, err = vault.Auth().Token().RenewSelf(creationTTL.(int))
		if err != nil {
			http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
			return
		}
		tokenTTL, err = secret.TokenTTL()
		if err != nil {
			http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
			return
		}
		span.SetAttributes(attribute.Stringer("auth.token_ttl_renewed", tokenTTL))

		sess.Options.MaxAge = int(tokenTTL.Seconds())
		if err := sess.Save(r, w); err != nil {
			http.Error(w, spanerr.RecordError(ctx, err).Error(), http.StatusInternalServerError)
			return
		}
	}

	vaultToken, err := secret.TokenID()
	if err != nil {
		http.Error(w, "no vault token", http.StatusForbidden)
		return
	}

	tokenMeta, err := secret.TokenMetadata()
	if err != nil {
		http.Error(w, "no token metadata", http.StatusForbidden)
		return
	}

	span.SetAttributes(attribute.String("auth.username", tokenMeta["username"]))

	w.Header().Set("X-Auth-Request-Token", vaultToken)
	w.Header().Set("X-Auth-Request-User", tokenMeta["username"])
	w.WriteHeader(http.StatusOK)
}
