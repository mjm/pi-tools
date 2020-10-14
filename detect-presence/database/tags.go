package database

import (
	"context"
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

func (c *Client) TagTrip(ctx context.Context, tripID string, tag Tag) error {
	if _, err := sq.Insert("trip_taggings").
		Columns("trip_id", "tag").
		Values(tripID, string(tag)).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return fmt.Errorf("tagging trip %s with tag %q: %w", tripID, tag, err)
	}

	return nil
}

func (c *Client) UntagTrip(ctx context.Context, tripID string, tag Tag) error {
	if _, err := sq.Delete("trip_taggings").
		Where(sq.Eq{
			"trip_id": tripID,
			"tag":     string(tag),
		}).
		RunWith(c.db).
		ExecContext(ctx); err != nil {
		return fmt.Errorf("untagging trip %s with tag %q: %w", tripID, tag, err)
	}

	return nil
}
