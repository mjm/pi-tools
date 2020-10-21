package migrate

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/mjm/pi-tools/pkg/migrate/embeddata"
)

func UpIfNeeded(db *sql.DB, files map[string][]byte) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	source, err := embeddata.WithFiles(files)
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
