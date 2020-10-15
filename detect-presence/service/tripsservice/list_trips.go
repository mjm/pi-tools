package tripsservice

import (
	"context"
	"fmt"
	"time"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) ListTrips(ctx context.Context, req *tripspb.ListTripsRequest) (*tripspb.ListTripsResponse, error) {
	res := &tripspb.ListTripsResponse{}

	trips, err := s.db.ListTrips(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing trips: %w", err)
	}

	for _, trip := range trips {
		t := &tripspb.Trip{
			Id: trip.ID,
		}

		t.LeftAt = trip.LeftAt.UTC().Format(time.RFC3339)
		if !trip.ReturnedAt.IsZero() {
			t.ReturnedAt = trip.ReturnedAt.UTC().Format(time.RFC3339)
		}

		for _, tag := range trip.Tags {
			t.Tags = append(t.Tags, string(tag))
		}

		res.Trips = append(res.Trips, t)
	}

	return res, nil
}
