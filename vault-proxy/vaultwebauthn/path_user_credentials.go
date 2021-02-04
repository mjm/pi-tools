package vaultwebauthn

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathUserCredentialsRequest(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("name") + "/credentials/request",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeNameString,
				Description: "Name of the user.",
			},
			"creation_response": {
				Type:        framework.TypeString,
				Description: "Credential creation response that should be sent to the browser, as JSON.",
			},
			"session_data": {
				Type:        framework.TypeString,
				Description: "Data that should be stored by the server during the authentication ceremony, as JSON.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathUserCredentialRequest,
			},
		},
	}
}

func pathUserCredentialsCreate(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("name") + "/credentials/create",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeNameString,
				Description: "Name of the user.",
			},
			"session_data": {
				Type:        framework.TypeString,
				Description: "Data that was stored by the server during the authentication ceremony, as JSON.",
				Required:    true,
			},
			"attestation_response": {
				Type:        framework.TypeString,
				Description: "Attestation response from the browser, as JSON",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathUserCredentialCreate,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUserCredentialCreate,
			},
		},
	}
}

func pathUserCredentialsList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("name") + "/credentials/?",
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeNameString,
				Description: "Name of the user.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathUserCredentialsList,
			},
		},
	}
}

func (b *backend) pathUserCredentialRequest(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	user, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return logical.ErrorResponse("no such user"), logical.ErrInvalidRequest
	}

	userWithCreds := &userWithCredentials{
		UserEntry: *user,
		name:      username,
	}
	wa, err := b.WebAuthn(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	creationResponse, sessionData, err := wa.BeginRegistration(userWithCreds)
	if err != nil {
		return nil, err
	}

	creationResponseBytes, err := json.Marshal(creationResponse)
	if err != nil {
		return nil, err
	}

	sessionDataBytes, err := json.Marshal(sessionData)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"creation_response": string(creationResponseBytes),
		"session_data":      string(sessionDataBytes),
	}
	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathUserCredentialCreate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

	attestationResponseRaw := d.Get("attestation_response").(string)
	creationData, err := protocol.ParseCredentialCreationResponseBody(strings.NewReader(attestationResponseRaw))
	if err != nil {
		return nil, err
	}

	userWithCreds := &userWithCredentials{
		UserEntry: *user,
		name:      username,
	}
	wa, err := b.WebAuthn(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	credential, err := wa.CreateCredential(userWithCreds, sessionData, creationData)
	if err != nil {
		return nil, err
	}

	credEntry := &CredentialEntry{
		ID:              credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Authenticator: AuthenticatorEntry{
			AAGUID:    credential.Authenticator.AAGUID,
			SignCount: credential.Authenticator.SignCount,
		},
	}

	if err := b.setCredential(ctx, req.Storage, username, credEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathUserCredentialsList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	creds, err := req.Storage.List(ctx, "user/"+username+"/credential/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(creds), nil
}

func (b *backend) userCredentials(ctx context.Context, s logical.Storage, username string) ([]*CredentialEntry, error) {
	credKeys, err := s.List(ctx, "user/"+username+"/credential/")
	if err != nil {
		return nil, err
	}

	var creds []*CredentialEntry
	for _, credKey := range credKeys {
		entry, err := s.Get(ctx, "user/"+username+"/credential/"+credKey)
		if err != nil {
			return nil, err
		}

		var result CredentialEntry
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, err
		}

		creds = append(creds, &result)
	}

	return creds, err
}

func (b *backend) setCredential(ctx context.Context, s logical.Storage, username string, credEntry *CredentialEntry) error {
	idStr := base64.URLEncoding.EncodeToString(credEntry.ID)
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("user/%s/credential/%s", username, idStr), credEntry)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

type CredentialEntry struct {
	ID              []byte             `json:"id"`
	PublicKey       []byte             `json:"public_key"`
	AttestationType string             `json:"attestation_type"`
	Authenticator   AuthenticatorEntry `json:"authenticator"`
}

type AuthenticatorEntry struct {
	AAGUID    []byte `json:"aaguid"`
	SignCount uint32 `json:"sign_count"`
}

type userWithCredentials struct {
	UserEntry
	name  string
	creds []*CredentialEntry
}

func (u *userWithCredentials) WebAuthnID() []byte {
	b, _ := uuid.ParseUUID(u.ID)
	return b
}

func (u *userWithCredentials) WebAuthnName() string {
	return u.name
}

func (u *userWithCredentials) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *userWithCredentials) WebAuthnIcon() string {
	return ""
}

func (u *userWithCredentials) WebAuthnCredentials() []webauthn.Credential {
	var creds []webauthn.Credential
	for _, savedCred := range u.creds {
		creds = append(creds, webauthn.Credential{
			ID:              savedCred.ID,
			PublicKey:       savedCred.PublicKey,
			AttestationType: savedCred.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:    savedCred.Authenticator.AAGUID,
				SignCount: savedCred.Authenticator.SignCount,
			},
		})
	}
	return creds
}
