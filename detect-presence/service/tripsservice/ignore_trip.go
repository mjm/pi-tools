package tripsservice

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) IgnoreTrip(ctx context.Context, req *tripspb.IgnoreTripRequest) (*tripspb.IgnoreTripResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID for trip to ignore")
	}

	if err := s.db.IgnoreTrip(ctx, req.Id); err != nil {
		return nil, err
	}

	return &tripspb.IgnoreTripResponse{}, nil
}
