package database

import (
	"context"
	"database/sql"

	"github.com/mjm/pi-tools/go-links/database/migrate"
	"github.com/mjm/pi-tools/pkg/migrate/postgres"
)

func (q *Queries) MigrateIfNeeded(_ context.Context) error {
	return postgres.UpIfNeeded(q.db.(*sql.DB), migrate.Data)
}
