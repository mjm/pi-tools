package database

import (
	"context"

	"github.com/mjm/pi-tools/pkg/migrate"
)

func (c *Client) MigrateIfNeeded(_ context.Context) error {
	return migrate.UpIfNeeded(c.db, nil)
}
