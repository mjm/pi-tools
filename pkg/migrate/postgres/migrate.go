package postgres

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/mjm/pi-tools/pkg/migrate/embeddata"
)

func UpIfNeeded(db *sql.DB, files map[string][]byte) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	source, err := embeddata.WithFiles(files)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-embed-data", source, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
