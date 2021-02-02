package authservice

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func (s *Server) StartOAuth(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	redirectURL := r.FormValue("redirect_uri")
	if redirectURL == "" {
		redirectURL = r.Header.Get("X-Auth-Request-Redirect")
	}
	if redirectURL == "" {
		redirectURL = "/"
	}

	url := s.OAuth.AuthCodeURL(state + ":" + redirectURL)
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  state,
		Path:   "/",
		Domain: s.CookieDomain,
	})
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
