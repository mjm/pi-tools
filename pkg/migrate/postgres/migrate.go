package postgres

import (
	"database/sql"
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	fssource "github.com/mjm/pi-tools/pkg/migrate/fs"
)

func UpIfNeeded(db *sql.DB, files fs.FS) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	source, err := fssource.WithFS(files)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("fs", source, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
