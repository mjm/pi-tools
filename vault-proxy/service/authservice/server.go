package authservice

import (
	vaultapi "github.com/hashicorp/vault/api"
	"golang.org/x/oauth2"
)

type Server struct {
	Vault        *vaultapi.Client
	OAuth        *oauth2.Config
	CookieDomain string
}

type Config struct {
	CookieDomain string
}

func New(vault *vaultapi.Client, oauth *oauth2.Config, cfg Config) (*Server, error) {
	return &Server{
		Vault:        vault,
		OAuth:        oauth,
		CookieDomain: cfg.CookieDomain,
	}, nil
}
