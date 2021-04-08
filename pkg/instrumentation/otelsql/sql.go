package otelsql

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultTracerName = "github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
)

type DB struct {
	*sql.DB
	cfg    *config
	tracer trace.Tracer
}

func NewDBWithTracing(db *sql.DB, opts ...Option) *DB {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &DB{
		DB:     db,
		cfg:    cfg,
		tracer: otel.Tracer(defaultTracerName),
	}
}

type config struct {
	baseAttrs []attribute.KeyValue
}

type Option func(*config)

func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(cfg *config) {
		cfg.baseAttrs = attrs
	}
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, span := db.tracer.Start(ctx, "ExecContext",
		trace.WithAttributes(db.cfg.baseAttrs...),
		trace.WithAttributes(semconv.DBStatementKey.String(query)))
	defer span.End()

	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return result, nil
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	ctx, span := db.tracer.Start(ctx, "PrepareContext",
		trace.WithAttributes(db.cfg.baseAttrs...),
		trace.WithAttributes(semconv.DBStatementKey.String(query)))
	defer span.End()

	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return stmt, nil
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, span := db.tracer.Start(ctx, "QueryContext",
		trace.WithAttributes(db.cfg.baseAttrs...),
		trace.WithAttributes(semconv.DBStatementKey.String(query)))
	defer span.End()

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return rows, nil
}

func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx, span := db.tracer.Start(ctx, "QueryRowContext",
		trace.WithAttributes(db.cfg.baseAttrs...),
		trace.WithAttributes(semconv.DBStatementKey.String(query)))
	defer span.End()

	// annoying: can't set the status on an error because we don't have access to it until the row
	// is scanned.
	return db.DB.QueryRowContext(ctx, query, args...)
}
