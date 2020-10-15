package tripsservice

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) GetTrip(ctx context.Context, req *tripspb.GetTripRequest) (*tripspb.GetTripResponse, error) {
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing trip ID")
	}

	trip, err := s.db.GetTrip(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

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

	return &tripspb.GetTripResponse{
		Trip: t,
	}, nil
}
