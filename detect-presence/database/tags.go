package database

import (
	"context"

	"github.com/google/uuid"
)

type UpdateTripTagsParams struct {
	TripID       uuid.UUID
	TagsToAdd    []string
	TagsToRemove []string
}

func (q *Queries) UpdateTripTags(ctx context.Context, args UpdateTripTagsParams) error {
	for _, tag := range args.TagsToAdd {
		if err := q.TagTrip(ctx, TagTripParams{
			TripID: args.TripID,
			Tag:    tag,
		}); err != nil {
			return err
		}
	}

	for _, tag := range args.TagsToRemove {
		if err := q.UntagTrip(ctx, UntagTripParams{
			TripID: args.TripID,
			Tag:    tag,
		}); err != nil {
			return err
		}
	}

	return nil
}
