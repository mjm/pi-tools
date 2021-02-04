package vaultwebauthn

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathUserAssertion(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "assertion",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeNameString,
				Description: "Name of the user.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathUserAssertion,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUserAssertion,
			},
		},
	}
}

func (b *backend) pathUserAssertion(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	user, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return logical.ErrorResponse("no such user"), logical.ErrInvalidRequest
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

	assertion, sessionData, err := wa.BeginLogin(userWithCreds)
	if err != nil {
		return nil, err
	}

	assertionBytes, err := json.Marshal(assertion)
	if err != nil {
		return nil, err
	}

	sessionDataBytes, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"assertion":    string(assertionBytes),
		"session_data": string(sessionDataBytes),
	}
	return &logical.Response{
		Data: data,
	}, nil
}
