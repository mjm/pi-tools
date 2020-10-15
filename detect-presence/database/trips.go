package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Trip struct {
	ID         string
	LeftAt     time.Time
	ReturnedAt time.Time
	Tags       []Tag
}

func (c *Client) BeginTrip(ctx context.Context, leftAt time.Time) (*Trip, error) {
	id := uuid.New().String()
	if _, err := sq.Insert("trips").
		Columns("id", "left_at").
		Values(id, leftAt.Unix()).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return nil, fmt.Errorf("beginning new trip: %w", err)
	}

	return &Trip{
		ID:     id,
		LeftAt: leftAt,
	}, nil
}

func (c *Client) EndTrip(ctx context.Context, id string, returnedAt time.Time) error {
	if _, err := sq.Update("trips").
		Set("returned_at", returnedAt.Unix()).
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Eq{"returned_at": nil},
		}).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return fmt.Errorf("ending trip: %w", err)
	}

	return nil
}

func (c *Client) GetCurrentTrip(ctx context.Context) (*Trip, error) {
	trip := new(Trip)
	var leftAtUnix int64
	if err := sq.Select("id", "left_at").
		From("trips").
		Where(sq.Eq{
			"ignored_at":  nil,
			"returned_at": nil,
		}).
		OrderBy("left_at DESC").
		Limit(1).
		RunWith(c.db).
		QueryRowContext(ctx).
		Scan(&trip.ID, &leftAtUnix); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("loading current trip: %w", err)
	}

	trip.LeftAt = time.Unix(leftAtUnix, 0)
	return trip, nil
}

func (c *Client) GetLastCompletedTrip(ctx context.Context) (*Trip, error) {
	trip := new(Trip)
	var leftAtUnix, returnedAtUnix int64
	if err := sq.Select("id", "left_at", "returned_at").
		From("trips").
		Where(sq.And{
			sq.Eq{"ignored_at": nil},
			sq.NotEq{"returned_at": nil},
		}).
		OrderBy("left_at DESC").
		Limit(1).
		RunWith(c.db).
		QueryRowContext(ctx).
		Scan(&trip.ID, &leftAtUnix, &returnedAtUnix); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("loading last completed trip: %w", err)
	}

	trip.LeftAt = time.Unix(leftAtUnix, 0)
	trip.ReturnedAt = time.Unix(returnedAtUnix, 0)
	return trip, nil
}

func (c *Client) ListTrips(ctx context.Context) ([]*Trip, error) {
	rows, err := sq.Select("id", "left_at", "returned_at", "group_concat(tag, '|') as tags").
		From("trips").
		LeftJoin("trip_taggings ON trips.id = trip_taggings.trip_id").
		Where(sq.Eq{"ignored_at": nil}).
		GroupBy("id").
		OrderBy("left_at DESC").
		Limit(30).
		RunWith(c.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing trips: %w", err)
	}

	var trips []*Trip
	for rows.Next() {
		trip, err := scanTrip(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning trip row: %w", err)
		}

		trips = append(trips, trip)
	}

	return trips, nil
}

func (c *Client) GetTrip(ctx context.Context, id string) (*Trip, error) {
	rows, err := sq.Select("id", "left_at", "returned_at", "group_concat(tag, '|') as tags").
		From("trips").
		LeftJoin("trip_taggings ON trips.id = trip_taggings.trip_id").
		Where(sq.Eq{"id": id}).
		GroupBy("id").
		RunWith(c.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting trip: %w", err)
	}

	if !rows.Next() {
		return nil, fmt.Errorf("no trip with ID %s", id)
	}

	trip, err := scanTrip(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning trip row: %w", err)
	}

	return trip, nil
}

func (c *Client) IgnoreTrip(ctx context.Context, id string) error {
	if _, err := sq.Update("trips").
		Set("ignored_at", time.Now().Unix()).
		Where(sq.Eq{"id": id}).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return fmt.Errorf("ignoring trip: %w", err)
	}

	return nil
}

func scanTrip(rows *sql.Rows) (*Trip, error) {
	trip := new(Trip)
	var leftAtUnix int64
	var returnedAtUnix *int64
	var tagsList *string

	if err := rows.Scan(&trip.ID, &leftAtUnix, &returnedAtUnix, &tagsList); err != nil {
		return nil, err
	}

	trip.LeftAt = time.Unix(leftAtUnix, 0)
	if returnedAtUnix != nil {
		trip.ReturnedAt = time.Unix(*returnedAtUnix, 0)
	}

	if tagsList != nil {
		tagStrings := strings.Split(*tagsList, "|")
		for _, tag := range tagStrings {
			trip.Tags = append(trip.Tags, Tag(tag))
		}
	}

	return trip, nil
}
