package storagetest

import (
	"context"
	"database/sql"

	"zombiezen.com/go/postgrestest"

	"github.com/mjm/pi-tools/pkg/migrate/postgres"
	"github.com/mjm/pi-tools/storage"
)

type wrappedDB struct {
	*sql.DB
}

func (wrappedDB) Check(context.Context) error {
	return nil
}

func NewDatabase(ctx context.Context, dbSrv *postgrestest.Server, migrations map[string][]byte) (storage.DB, error) {
	db, err := dbSrv.NewDatabase(ctx)
	if err != nil {
		return nil, err
	}

	if err := postgres.UpIfNeeded(db, migrations); err != nil {
		return nil, err
	}

	return &wrappedDB{DB: db}, nil
}
