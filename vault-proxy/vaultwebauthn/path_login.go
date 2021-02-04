package vaultwebauthn

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	for _, credEntry := range creds {
		if !bytes.Equal(credEntry.ID, usedCred.ID) {
			continue
		}

		credEntry.Authenticator.SignCount = usedCred.Authenticator.SignCount
		if err := b.setCredential(ctx, req.Storage, username, credEntry); err != nil {
			return nil, err
		}
		break
	}

	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	auth := &logical.Auth{
		InternalData: map[string]interface{}{
			"credential_id": base64.URLEncoding.EncodeToString(usedCred.ID),
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

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Auth == nil {
		return nil, fmt.Errorf("request auth was nil")
	}

	credentialIDRaw, ok := req.Auth.InternalData["credential_id"]
	if !ok {
		return nil, fmt.Errorf("no credential ID associated with token")
	}
	credentialID, err := base64.URLEncoding.DecodeString(credentialIDRaw.(string))
	if err != nil {
		return nil, err
	}

	// verify that this credential is still present for the user
	username, ok := req.Auth.Metadata["username"]
	if !ok {
		return nil, fmt.Errorf("no username associated with token")
	}
	creds, err := b.userCredentials(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}

	var foundCred *CredentialEntry
	for _, cred := range creds {
		if !bytes.Equal(cred.ID, credentialID) {
			continue
		}

		foundCred = cred
		break
	}

	if foundCred == nil {
		return nil, fmt.Errorf("credential used to obtain token is no longer registered to user")
	}

	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.Period = cfg.TokenPeriod
	resp.Auth.TTL = cfg.TokenTTL
	resp.Auth.MaxTTL = cfg.TokenMaxTTL

	return resp, nil
}
