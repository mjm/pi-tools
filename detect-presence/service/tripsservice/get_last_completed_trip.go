package tripsservice

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) GetLastCompletedTrip(ctx context.Context, _ *tripspb.GetLastCompletedTripRequest) (*tripspb.GetLastCompletedTripResponse, error) {
	span := trace.SpanFromContext(ctx)

	trip, err := s.q.GetLastCompletedTrip(ctx)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.String("trip.id", trip.ID.String()))

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
