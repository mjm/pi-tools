package vaultwebauthn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathUsersList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/?",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathUserList,
			},
		},
	}
}

func pathUsers(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeNameString,
				Description: "Name of the user.",
			},
			"id": {
				Type:        framework.TypeLowerCaseString,
				Description: "ID of the user, auto-generated.",
			},
			"display_name": {
				Type:        framework.TypeString,
				Description: "Display name of the user.",
			},
			// TODO optional icon
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathUserRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathUserWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathUserWrite,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathUserDelete,
			},
		},
		ExistenceCheck: b.userExistenceCheck,
	}
	tokenutil.AddTokenFields(p.Fields)
	return p
}

func (b *backend) userExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	userEntry, err := b.user(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return false, err
	}

	return userEntry != nil, nil
}

func (b *backend) pathUserList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	users, err := req.Storage.List(ctx, "user/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(users), nil
}

func (b *backend) pathUserDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "user/"+strings.ToLower(d.Get("name").(string)))
	if err != nil {
		return nil, err
	}

	// TODO also delete the user's credentials

	return nil, nil
}

func (b *backend) pathUserRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	user, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	data := map[string]interface{}{}
	user.PopulateTokenData(data)

	data["name"] = username
	data["id"] = user.ID
	data["display_name"] = user.DisplayName

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathUserWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := strings.ToLower(d.Get("name").(string))
	userEntry, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}

	// Due to existence check, user will only be nil if it's a create operation
	if userEntry == nil {
		userUUID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		userEntry = &UserEntry{
			ID: userUUID,
		}
	}

	if err := userEntry.ParseTokenFields(req, d); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	if displayNameRaw, ok := d.GetOk("display_name"); ok {
		userEntry.DisplayName = displayNameRaw.(string)
	}

	return nil, b.setUser(ctx, req.Storage, username, userEntry)
}

func (b *backend) user(ctx context.Context, s logical.Storage, username string) (*UserEntry, error) {
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	entry, err := s.Get(ctx, "user/"+strings.ToLower(username))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result UserEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) setUser(ctx context.Context, s logical.Storage, username string, userEntry *UserEntry) error {
	entry, err := logical.StorageEntryJSON("user/"+username, userEntry)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

type UserEntry struct {
	tokenutil.TokenParams

	ID          string `json:"id" structs:"id" mapstructure:"id"`
	DisplayName string `json:"display_name" structs:"display_name" mapstructure:"display_name"`
}
