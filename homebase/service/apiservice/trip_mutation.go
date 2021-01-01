package apiservice

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (r *Resolver) IgnoreTrip(ctx context.Context, args struct {
	Input struct {
		ID graphql.ID
	}
}) (
	resp struct {
		IgnoredTripID graphql.ID
	},
	err error,
) {
	var id string
	if err = relay.UnmarshalSpec(args.Input.ID, &id); err != nil {
		return
	}

	if _, err = r.tripsClient.IgnoreTrip(ctx, &tripspb.IgnoreTripRequest{Id: id}); err != nil {
		return
	}

	resp.IgnoredTripID = args.Input.ID
	return
}

func (r *Resolver) UpdateTripTags(ctx context.Context, args struct {
	Input struct {
		TripID       graphql.ID
		TagsToAdd    []string
		TagsToRemove []string
	}
}) (
	resp struct {
		Trip *Trip
	},
	err error,
) {
	var id string
	if err = relay.UnmarshalSpec(args.Input.TripID, &id); err != nil {
		return
	}

	if _, err = r.tripsClient.UpdateTripTags(ctx, &tripspb.UpdateTripTagsRequest{
		TripId:       id,
		TagsToAdd:    args.Input.TagsToAdd,
		TagsToRemove: args.Input.TagsToRemove,
	}); err != nil {
		return
	}

	resp.Trip, err = r.Trip(ctx, struct{ ID graphql.ID }{ID: args.Input.TripID}) // oof
	return
}