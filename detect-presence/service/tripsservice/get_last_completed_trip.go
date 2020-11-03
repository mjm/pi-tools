package tripsservice

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) GetLastCompletedTrip(ctx context.Context, _ *tripspb.GetLastCompletedTripRequest) (*tripspb.GetLastCompletedTripResponse, error) {
	span := trace.SpanFromContext(ctx)

	trip, err := s.q.GetLastCompletedTrip(ctx)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(label.String("trip.id", trip.ID.String()))

	t := &tripspb.Trip{
		Id:     trip.ID.String(),
		LeftAt: trip.LeftAt.UTC().Format(time.RFC3339),
	}

	if trip.ReturnedAt.Valid {
		t.ReturnedAt = trip.ReturnedAt.Time.UTC().Format(time.RFC3339)
	}

	return &tripspb.GetLastCompletedTripResponse{
		Trip: t,
	}, nil
}
