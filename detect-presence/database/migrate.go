package database

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	migrations "github.com/mjm/pi-tools/detect-presence/database/migrate"
)

func (c *Client) MigrateIfNeeded(_ context.Context) error {
	driver, err := sqlite3.WithInstance(c.db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	source, err := newEmbeddedMigrationsSource(migrations.Data)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-embed-data", source, "sqlite3", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
