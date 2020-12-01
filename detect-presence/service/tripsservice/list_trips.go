package tripsservice

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) ListTrips(ctx context.Context, req *tripspb.ListTripsRequest) (*tripspb.ListTripsResponse, error) {
	span := trace.SpanFromContext(ctx)

	res := &tripspb.ListTripsResponse{}

	var limit int32 = 30
	if req.GetLimit() > 0 && req.GetLimit() < 100 {
		limit = req.GetLimit()
	}

	span.SetAttributes(label.Int32("limit", limit))
	trips, err := s.q.ListTrips(ctx, limit)
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
