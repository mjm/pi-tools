package authservice

import (
	"net/http"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

func (s *Server) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	authzHeader := r.Header.Get("Authorization")
	if authzHeader != "" {
		if !strings.HasPrefix(authzHeader, "Bearer ") {
			http.Error(w, "invalid Authorization header value", http.StatusForbidden)
			return
		}

		token := authzHeader[7:]
		s.handleGitHubToken(w, token)
		return
	}

	vaultTokenHeader := r.Header.Get("X-Vault-Token")
	if vaultTokenHeader != "" {
		s.handleVaultToken(w, vaultTokenHeader)
		return
	}

	vaultProxyCookie, err := r.Cookie("vault_proxy")
	if err == nil {
		s.handleVaultToken(w, vaultProxyCookie.Value)
		return
	}

	http.Error(w, "no credentials present", http.StatusUnauthorized)
}

func (s *Server) handleGitHubToken(w http.ResponseWriter, token string) {
	secret, err := s.Vault.Logical().Write("auth/github-batch/login", map[string]interface{}{
		"token": token,
	})
	if err != nil {
		http.Error(w, "error logging in to GitHub", http.StatusForbidden)
		return
	}

	s.writeAuthResponse(w, secret)
}

func (s *Server) handleVaultToken(w http.ResponseWriter, token string) {
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
	http.SetCookie(w, &http.Cookie{
		Name:   "vault_proxy",
		Value:  vaultToken,
		Path:   "/",
		Domain: s.CookieDomain,
	})
	w.WriteHeader(http.StatusOK)
}
