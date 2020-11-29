package storage

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/etherlabsio/healthcheck"
	"go.opentelemetry.io/otel/semconv"

	"github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
	"github.com/mjm/pi-tools/pkg/migrate/postgres"
)

var dbDSN *string

// SetDefaultDBName configures the current process to support connecting to a PostgreSQL database. The name should be
// the name of the database used for development. This should be called before calling flag.Parse(), since it creates
// the -db CLI flag for specifying the database DSN in production.
func SetDefaultDBName(name string) {
	dbDSN = flag.String(
		"db",
		fmt.Sprintf("dbname=%s sslmode=disable", name),
		"Connection string for connecting to PostgreSQL database")
}

// DB is an interface that includes a subset of the methods on sql.DB. Code that works with DB connections
// should use this interface type rather than *sql.DB directly so that we can substitute an implementation
// that supports tracing.
type DB interface {
	healthcheck.Checker
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type wrappedDB struct {
	*otelsql.DB
}

// OpenDB opens the PostgreSQL database for the current application. Panics if SetDefaultDBName has not been called.
// Pass in a map of migration file contents that were embedded into the binary at build time. These migrations will be
// run immediately after opening the database to ensure the schema is up-to-date. Returns a database that is
// instrumented for tracing.
func OpenDB(migrations map[string][]byte) (DB, error) {
	if dbDSN == nil {
		log.Panicf("no default database name configured")
	}

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

	return &wrappedDB{DB: db}, nil
}

// MustOpenDB calls OpenDB, but panics if there is an error opening the database or running migrations. Use this in main
// functions, when panicking was going to be your solution to an error anyway.
func MustOpenDB(migrations map[string][]byte) DB {
	db, err := OpenDB(migrations)
	if err != nil {
		log.Panicf("Error setting up storage: %v", err)
	}
	return db
}
