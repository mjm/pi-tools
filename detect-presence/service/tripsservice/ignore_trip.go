package tripsservice

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) IgnoreTrip(ctx context.Context, req *tripspb.IgnoreTripRequest) (*tripspb.IgnoreTripResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.String("trip.id", req.GetId()))

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID for trip to ignore")
	}

	tripID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID for trip ID: %s", err)
	}

	n, err := s.q.IgnoreTrip(ctx, tripID)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, status.Errorf(codes.NotFound, "no trip found with ID %s", tripID)
	}

	return &tripspb.IgnoreTripResponse{}, nil
}
