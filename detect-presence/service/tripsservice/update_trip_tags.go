package tripsservice

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/detect-presence/database"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) UpdateTripTags(ctx context.Context, req *tripspb.UpdateTripTagsRequest) (*tripspb.UpdateTripTagsResponse, error) {
	if req.GetTripId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID for trip to tag")
	}

	var tagsToAdd, tagsToRemove []database.Tag
	for _, tag := range req.GetTagsToAdd() {
		tagsToAdd = append(tagsToAdd, database.Tag(tag))
	}
	for _, tag := range req.GetTagsToRemove() {
		tagsToRemove = append(tagsToRemove, database.Tag(tag))
	}

	if err := s.db.UpdateTripTags(ctx, req.GetTripId(), tagsToAdd, tagsToRemove); err != nil {
		return nil, fmt.Errorf("updating trip tags: %w", err)
	}

	return &tripspb.UpdateTripTagsResponse{}, nil
}
