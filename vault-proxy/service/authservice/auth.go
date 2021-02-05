package authservice

import (
	"net/http"

	"github.com/gorilla/sessions"
	vaultapi "github.com/hashicorp/vault/api"
)

func (s *Server) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	vaultTokenHeader := r.Header.Get("X-Vault-Token")
	if vaultTokenHeader != "" {
		s.handleVaultToken(w, vaultTokenHeader, nil)
		return
	}

	sess, err := s.Store.Get(r, "vault-token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if vaultTokenRaw, ok := sess.Values["token"]; ok {
		s.handleVaultToken(w, vaultTokenRaw.(string), sess)
		return
	}

	http.Error(w, "no credentials present", http.StatusUnauthorized)
}

func (s *Server) handleVaultToken(w http.ResponseWriter, token string, sess *sessions.Session) {
	vault, err := s.Vault.Clone()
	if err != nil {
		http.Error(w, "failed to clone vault client", http.StatusInternalServerError)
		return
	}

	vault.SetToken(token)
	secret, err := vault.Auth().Token().LookupSelf()
	if err != nil {
		http.Error(w, "invalid vault token", http.StatusForbidden)
		return
	}

	// TODO renew token if needed

	s.writeAuthResponse(w, secret)
}

func (s *Server) writeAuthResponse(w http.ResponseWriter, secret *vaultapi.Secret) {
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

	w.Header().Set("X-Auth-Request-Token", vaultToken)
	w.Header().Set("X-Auth-Request-User", tokenMeta["username"])
	w.WriteHeader(http.StatusOK)
}
