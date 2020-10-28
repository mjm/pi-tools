// Code generated by sqlc. DO NOT EDIT.

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Trip struct {
	ID         uuid.UUID
	LeftAt     time.Time
	ReturnedAt sql.NullTime
	IgnoredAt  sql.NullTime
}

type TripTagging struct {
	TripID uuid.UUID
	Tag    string
}