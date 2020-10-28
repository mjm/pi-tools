// Code generated by sqlc. DO NOT EDIT.
// source: tags.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const tagTrip = `-- name: TagTrip :exec
INSERT INTO trip_taggings (trip_id, tag)
VALUES ($1, $2)
`

type TagTripParams struct {
	TripID uuid.UUID
	Tag    string
}

func (q *Queries) TagTrip(ctx context.Context, arg TagTripParams) error {
	_, err := q.db.ExecContext(ctx, tagTrip, arg.TripID, arg.Tag)
	return err
}

const untagTrip = `-- name: UntagTrip :exec
DELETE FROM trip_taggings
WHERE trip_id = $1
AND tag = $2
`

type UntagTripParams struct {
	TripID uuid.UUID
	Tag    string
}

func (q *Queries) UntagTrip(ctx context.Context, arg UntagTripParams) error {
	_, err := q.db.ExecContext(ctx, untagTrip, arg.TripID, arg.Tag)
	return err
}