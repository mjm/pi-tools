package apiservice

import (
	"time"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

type TripConnection struct {
	res *tripspb.ListTripsResponse
}

func (tc *TripConnection) Edges() []TripEdge {
	var edges []TripEdge
	for _, t := range tc.res.GetTrips() {
		edges = append(edges, TripEdge{Trip: t})
	}
	return edges
}

func (tc *TripConnection) PageInfo() PageInfo {
	return PageInfo{}
}

func (tc *TripConnection) TotalCount() int32 {
	return int32(len(tc.res.GetTrips()))
}

type TripEdge struct {
	*tripspb.Trip
}

func (e TripEdge) Node() *Trip {
	return &Trip{Trip: e.Trip}
}

func (e TripEdge) Cursor() Cursor {
	return Cursor(e.GetId())
}

type Trip struct {
	*tripspb.Trip
}

func (t *Trip) ID() graphql.ID {
	return relay.MarshalID("trip", t.GetId())
}

func (t *Trip) LeftAt() (graphql.Time, error) {
	t2, err := time.Parse(time.RFC3339, t.GetLeftAt())
	if err != nil {
		return graphql.Time{}, err
	}
	return graphql.Time{Time: t2}, nil
}

func (t *Trip) ReturnedAt() (*graphql.Time, error) {
	if t.GetReturnedAt() == "" {
		return nil, nil
	}
	t2, err := time.Parse(time.RFC3339, t.GetReturnedAt())
	if err != nil {
		return nil, err
	}
	return &graphql.Time{Time: t2}, nil
}

func (t *Trip) Tags() []string {
	return t.GetTags()
}
