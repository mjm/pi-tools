package database

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Link struct {
	ID             string
	ShortURL       string
	DestinationURL string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateLinkParams struct {
	ShortURL       string
	DestinationURL string
	Description    string
}

func (c *Client) CreateLink(ctx context.Context, params CreateLinkParams) (*Link, error) {
	id := uuid.New().String()
	now := time.Now()

	if _, err := sq.Insert("links").
		Columns("id", "short_url", "destination_url", "description", "created_at", "updated_at").
		Values(id, params.ShortURL, params.DestinationURL, params.Description, now.Unix(), now.Unix()).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return nil, fmt.Errorf("creating link: %w", err)
	}

	return &Link{
		ID:             id,
		ShortURL:       params.ShortURL,
		DestinationURL: params.DestinationURL,
		Description:    params.Description,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
