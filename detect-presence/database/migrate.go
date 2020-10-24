package database

import (
	"context"

	"github.com/mjm/pi-tools/detect-presence/database/migrate"
	"github.com/mjm/pi-tools/pkg/migrate/sqlite3"
)

func (c *Client) MigrateIfNeeded(_ context.Context) error {
	return sqlite3.UpIfNeeded(c.db, migrate.Data)
}
