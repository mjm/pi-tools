package apiservice

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

type Resolver struct {
	tripsClient tripspb.TripsServiceClient
}

func (r *Resolver) Viewer() *Resolver {
	return r
}

func (r *Resolver) Node(ctx context.Context, args struct{ ID graphql.ID }) (*Node, error) {
	return nil, nil
}

func (r *Resolver) Trips(ctx context.Context, args struct {
	First *int32
	After *string
}) (*TripConnection, error) {
	// TODO actually support paging

	var limit int32 = 30
	if args.First != nil {
		limit = *args.First
	}
	res, err := r.tripsClient.ListTrips(ctx, &tripspb.ListTripsRequest{Limit: limit})
	if err != nil {
		return nil, err
	}
	return &TripConnection{res: res}, nil
}

func (r *Resolver) Trip(ctx context.Context, args struct {
	ID graphql.ID
}) (*Trip, error) {
	var id string
	if err := relay.UnmarshalSpec(args.ID, &id); err != nil {
		return nil, err
	}

	res, err := r.tripsClient.GetTrip(ctx, &tripspb.GetTripRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &Trip{Trip: res.GetTrip()}, nil
}
