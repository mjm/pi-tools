package vaultwebauthn

import (
	"context"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// Factory returns a configured instance of the backend.
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		BackendType: logical.TypeCredential,
		AuthRenew:   b.pathLoginRenew,
		Help:        "",
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"assertion",
				"login",
			},
		},
		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathLogin(&b),
			pathUsersList(&b),
			pathUsers(&b),
			pathUserCredentialsList(&b),
			pathUserCredentialsRequest(&b),
			pathUserCredentialsCreate(&b),
			pathUserAssertion(&b),
		}),
	}

	return &b
}

func (b *backend) WebAuthn(ctx context.Context, s logical.Storage) (*webauthn.WebAuthn, error) {
	cfg, err := b.Config(ctx, s)
	if err != nil {
		return nil, err
	}

	wa, err := webauthn.New(&webauthn.Config{
		RPDisplayName: cfg.DisplayName,
		RPID:          cfg.ID,
		RPOrigin:      cfg.Origin,
	})
	if err != nil {
		return nil, err
	}
	return wa, nil
}
