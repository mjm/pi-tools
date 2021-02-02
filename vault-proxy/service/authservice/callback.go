package authservice

import (
	"net/http"
	"strings"
)

func (s *Server) HandleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	state := strings.SplitN(r.FormValue("state"), ":", 2)
	if len(state) < 2 {
		http.Error(w, "invalid state", http.StatusForbidden)
		return
	}
	nonce := state[0]
	redirectURL := state[1]

	savedState, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nonce != savedState.Value {
		http.Error(w, "nonce mismatch", http.StatusForbidden)
		return
	}

	code := r.FormValue("code")
	token, err := s.OAuth.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	secret, err := s.Vault.Logical().Write("auth/github/login", map[string]interface{}{
		"token": token.AccessToken,
	})
	if err != nil {
		http.Error(w, "error logging in to GitHub: "+err.Error(), http.StatusForbidden)
		return
	}

	vaultToken, err := secret.TokenID()
	if err != nil {
		http.Error(w, "no vault token", http.StatusForbidden)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "vault_proxy",
		Value:  vaultToken,
		Path:   "/",
		Domain: s.CookieDomain,
	})
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
