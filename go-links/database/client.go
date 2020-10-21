package database

import (
	"database/sql"
	"fmt"
)

type Client struct {
	db *sql.DB
}

func Open(dsn string) (*Client, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	return &Client{db: db}, nil
}
