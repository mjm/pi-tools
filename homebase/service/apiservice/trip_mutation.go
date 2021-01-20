package apiservice

import (
	"context"
	"time"

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

func (r *Resolver) RecordTrips(ctx context.Context, args struct {
	Input struct {
		Trips []struct {
			ID         string
			LeftAt     graphql.Time
			ReturnedAt graphql.Time
		}
	}
}) (
	resp struct {
		RecordedTrips []*Trip
		Failures      []TripRecordingFailure
	},
	err error,
) {
	req := &tripspb.RecordTripsRequest{}
	for _, t := range args.Input.Trips {
		req.Trips = append(req.Trips, &tripspb.Trip{
			Id:         t.ID,
			LeftAt:     t.LeftAt.Format(time.RFC3339),
			ReturnedAt: t.ReturnedAt.Format(time.RFC3339),
		})
	}

	var res *tripspb.RecordTripsResponse
	res, err = r.tripsClient.RecordTrips(ctx, req)
	if err != nil {
		return
	}

	failedIDs := map[string]struct{}{}
	for _, failure := range res.GetFailures() {
		failedIDs[failure.GetTripId()] = struct{}{}
		resp.Failures = append(resp.Failures, TripRecordingFailure{
			TripID:  failure.GetTripId(),
			Message: failure.GetMessage(),
		})
	}

	for _, t := range req.Trips {
		if _, ok := failedIDs[t.GetId()]; ok {
			continue
		}

		resp.RecordedTrips = append(resp.RecordedTrips, &Trip{t})
	}
	return
}

type TripRecordingFailure struct {
	TripID  string
	Message string
}
