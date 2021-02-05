package authservice

import (
	"encoding/base64"

	"github.com/gorilla/sessions"
	vaultapi "github.com/hashicorp/vault/api"
)

type Server struct {
	Vault        *vaultapi.Client
	Store        sessions.Store
	AuthPath     string
	CookieDomain string
}

type Config struct {
	AuthPath     string
	CookieDomain string
	CookieKey    string
}

func New(vault *vaultapi.Client, cfg Config) (*Server, error) {
	key, err := base64.StdEncoding.DecodeString(cfg.CookieKey)
	if err != nil {
		return nil, err
	}

	store := sessions.NewCookieStore(key)
	return &Server{
		Vault:        vault,
		Store:        store,
		AuthPath:     cfg.AuthPath,
		CookieDomain: cfg.CookieDomain,
	}, nil
}
