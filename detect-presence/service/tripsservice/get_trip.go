package tripsservice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) GetTrip(ctx context.Context, req *tripspb.GetTripRequest) (*tripspb.GetTripResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.String("trip.id", req.GetId()))

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing trip ID")
	}

	tripID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID for trip ID: %s", err)
	}

	trip, err := s.q.GetTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}

	t := &tripspb.Trip{
		Id:     trip.ID.String(),
		LeftAt: trip.LeftAt.UTC().Format(time.RFC3339),
		Tags:   trip.Tags,
	}

	if trip.ReturnedAt.Valid {
		t.ReturnedAt = trip.ReturnedAt.Time.UTC().Format(time.RFC3339)
	}

	return &tripspb.GetTripResponse{
		Trip: t,
	}, nil
}
