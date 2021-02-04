package vaultwebauthn

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeNameString,
				Description: "Name of the user.",
				Required:    true,
			},
			"session_data": {
				Type:        framework.TypeString,
				Description: "Data that was stored by the server during the authentication ceremony, as JSON.",
				Required:    true,
			},
			"assertion_response": {
				Type:        framework.TypeString,
				Description: "Assertion response from the browser, as JSON",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathLogin,
			},
		},
	}
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	user, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return logical.ErrorResponse("no such user"), logical.ErrInvalidRequest
	}

	sessionDataRaw := d.Get("session_data").(string)

	var sessionData webauthn.SessionData
	if err := json.Unmarshal([]byte(sessionDataRaw), &sessionData); err != nil {
		return nil, err
	}

	assertionResponseRaw := d.Get("assertion_response").(string)
	assertionData, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(assertionResponseRaw))
	if err != nil {
		return nil, err
	}

	creds, err := b.userCredentials(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}

	userWithCreds := &userWithCredentials{
		UserEntry: *user,
		name:      username,
		creds:     creds,
	}
	wa, err := b.WebAuthn(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	usedCred, err := wa.ValidateLogin(userWithCreds, sessionData, assertionData)
	if err != nil {
		b.Logger().Error(err.Error())
		return nil, logical.ErrPermissionDenied
	}

	// TODO update sign count on the corresponding credential entry
	b.Logger().Info("validated login", "cred_id", usedCred.ID)

	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	auth := &logical.Auth{
		InternalData: map[string]interface{}{
			"credential_id": usedCred.ID,
		},
		Metadata: map[string]string{
			"username": username,
		},
		DisplayName: username,
		Alias: &logical.Alias{
			Name: username,
		},
	}
	cfg.PopulateTokenAuth(auth)
	auth.Policies = append(auth.Policies, user.TokenPolicies...)

	return &logical.Response{
		Auth: auth,
	}, nil
}
