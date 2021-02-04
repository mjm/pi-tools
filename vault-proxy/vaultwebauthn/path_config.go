package vaultwebauthn

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfig(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"display_name": {
				Type:        framework.TypeString,
				Description: "The display name for the site",
			},
			"id": {
				Type:        framework.TypeString,
				Description: "The FQDN for the site",
			},
			"origin": {
				Type:        framework.TypeString,
				Description: "The URL for WebAuthn requests",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
			},
		},
	}
	tokenutil.AddTokenFields(p.Fields)
	return p
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	c, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if c == nil {
		c = &config{}
	}

	if displayNameRaw, ok := data.GetOk("display_name"); ok {
		c.DisplayName = displayNameRaw.(string)
	}

	if idRaw, ok := data.GetOk("id"); ok {
		c.ID = idRaw.(string)
	}

	if originRaw, ok := data.GetOk("origin"); ok {
		c.Origin = originRaw.(string)
	}

	if err := c.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	entry, err := logical.StorageEntryJSON("config", c)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	d := map[string]interface{}{
		"display_name": config.DisplayName,
		"id":           config.ID,
		"origin":       config.Origin,
	}
	config.PopulateTokenData(d)

	return &logical.Response{
		Data: d,
	}, nil
}

// Config returns the configuration for this backend.
func (b *backend) Config(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, errwrap.Wrapf("error reading configuration: {{err}}", err)
	}

	return &result, nil
}

type config struct {
	tokenutil.TokenParams

	DisplayName string `json:"display_name" structs:"display_name" mapstructure:"display_name"`
	ID          string `json:"id" structs:"id" mapstructure:"id"`
	Origin      string `json:"origin" structs:"origin" mapstructure:"origin"`
}
