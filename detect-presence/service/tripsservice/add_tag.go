package tripsservice

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/detect-presence/database"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) TagTrip(ctx context.Context, req *tripspb.TagTripRequest) (*tripspb.TagTripResponse, error) {
	if req.GetTripId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID for trip to tag")
	}

	if req.GetTag() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tag")
	}

	if err := s.db.TagTrip(ctx, req.GetTripId(), database.Tag(req.GetTag())); err != nil {
		return nil, err
	}

	return &tripspb.TagTripResponse{}, nil
}
