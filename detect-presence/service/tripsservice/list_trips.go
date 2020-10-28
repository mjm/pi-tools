package tripsservice

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) ListTrips(ctx context.Context, req *tripspb.ListTripsRequest) (*tripspb.ListTripsResponse, error) {
	span := trace.SpanFromContext(ctx)

	res := &tripspb.ListTripsResponse{}

	trips, err := s.q.ListTrips(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "listing trips: %w", err)
	}

	span.SetAttributes(label.Int("trip.count", len(trips)))

	for _, trip := range trips {
		t := &tripspb.Trip{
			Id:     trip.ID.String(),
			LeftAt: trip.LeftAt.UTC().Format(time.RFC3339),
			Tags:   trip.Tags,
		}

		if trip.ReturnedAt.Valid {
			t.ReturnedAt = trip.ReturnedAt.Time.UTC().Format(time.RFC3339)
		}

		res.Trips = append(res.Trips, t)
	}

	return res, nil
}
