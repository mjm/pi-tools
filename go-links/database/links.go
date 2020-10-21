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
	if _, err := sq.Insert("links").
		Columns("id", "short_url", "destination_url", "description").
		Values(id, params.ShortURL, params.DestinationURL, params.Description).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return nil, fmt.Errorf("creating link: %w", err)
	}

	return &Link{
		ID:             id,
		ShortURL:       params.ShortURL,
		DestinationURL: params.DestinationURL,
		Description:    params.Description,
	}, nil
}
