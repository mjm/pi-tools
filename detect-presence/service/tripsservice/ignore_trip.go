package tripsservice

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) IgnoreTrip(ctx context.Context, req *tripspb.IgnoreTripRequest) (*tripspb.IgnoreTripResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.String("trip.id", req.GetId()))

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID for trip to ignore")
	}

	if err := s.db.IgnoreTrip(ctx, req.Id); err != nil {
		return nil, err
	}

	return &tripspb.IgnoreTripResponse{}, nil
}
