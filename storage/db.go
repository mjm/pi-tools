package storage

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/semconv"

	"github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
	"github.com/mjm/pi-tools/pkg/migrate/postgres"
)

var dbDSN *string

func SetDefaultDBName(name string) {
	dbDSN = flag.String(
		"db",
		fmt.Sprintf("dbname=%s sslmode=disable", name),
		"Connection string for connecting to PostgreSQL database")
}

type DB interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func OpenDB(migrations map[string][]byte) (DB, error) {
	sqlDB, err := sql.Open("postgres", *dbDSN)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	db := otelsql.NewDBWithTracing(sqlDB,
		otelsql.WithAttributes(
			semconv.DBSystemPostgres,
			// assuming this is safe to include since it was on the command-line.
			// passwords should come from a file or environment variable.
			semconv.DBConnectionStringKey.String(*dbDSN)))

	if err := postgres.UpIfNeeded(sqlDB, migrations); err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	return db, nil
}

func MustOpenDB(migrations map[string][]byte) DB {
	db, err := OpenDB(migrations)
	if err != nil {
		log.Panicf("Error setting up storage: %v", err)
	}
	return db
}
