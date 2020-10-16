package database

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type Tag string

func (c *Client) ListTags(ctx context.Context) ([]Tag, error) {
	rows, err := sq.Select("tag").
		From("trip_taggings").
		GroupBy("tag").
		OrderBy("count(trip_id) desc").
		RunWith(c.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing tags: %w", err)
	}

	var tags []Tag
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, Tag(tag))
	}
	return tags, nil
}

func (c *Client) tagTrip(ctx context.Context, tx *sql.Tx, tripID string, tag Tag) error {
	if _, err := sq.Insert("trip_taggings").
		Columns("trip_id", "tag").
		Values(tripID, string(tag)).
		RunWith(tx).
		ExecContext(ctx); err != nil {
		return fmt.Errorf("tagging trip %s with tag %q: %w", tripID, tag, err)
	}

	return nil
}

func (c *Client) untagTrip(ctx context.Context, tx *sql.Tx, tripID string, tag Tag) error {
	if _, err := sq.Delete("trip_taggings").
		Where(sq.Eq{
			"trip_id": tripID,
			"tag":     string(tag),
		}).
		RunWith(tx).
		ExecContext(ctx); err != nil {
		return fmt.Errorf("untagging trip %s with tag %q: %w", tripID, tag, err)
	}

	return nil
}

func (c *Client) UpdateTripTags(ctx context.Context, tripID string, tagsToAdd []Tag, tagsToRemove []Tag) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tag := range tagsToAdd {
		if err := c.tagTrip(ctx, tx, tripID, tag); err != nil {
			return err
		}
	}

	for _, tag := range tagsToRemove {
		if err := c.untagTrip(ctx, tx, tripID, tag); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing tag updates: %w", err)
	}

	return nil
}
