package tripsservice

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/detect-presence/database"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

func (s *Server) UpdateTripTags(ctx context.Context, req *tripspb.UpdateTripTagsRequest) (*tripspb.UpdateTripTagsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		label.String("trip.id", req.GetTripId()),
		label.Array("trip.tags.added", req.GetTagsToAdd()),
		label.Array("trip.tags.removed", req.GetTagsToRemove()))

	if req.GetTripId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID for trip to tag")
	}

	tripID, err := uuid.Parse(req.GetTripId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID for trip ID: %s", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "starting transaction: %s", err)
	}
	defer tx.Rollback()

	q := s.q.WithTx(tx)
	if err := q.UpdateTripTags(ctx, database.UpdateTripTagsParams{
		TripID:       tripID,
		TagsToAdd:    req.GetTagsToAdd(),
		TagsToRemove: req.GetTagsToRemove(),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "updating trip tags: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "updating trip tags: %w", err)
	}

	return &tripspb.UpdateTripTagsResponse{}, nil
}
