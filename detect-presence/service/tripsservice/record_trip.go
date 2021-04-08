package tripsservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/detect-presence/database"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
)

func (s *Server) RecordTrips(ctx context.Context, req *tripspb.RecordTripsRequest) (*tripspb.RecordTripsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("trip.count", len(req.GetTrips())))

	if len(req.GetTrips()) == 0 {
		return &tripspb.RecordTripsResponse{}, nil
	}

	var failures []*tripspb.RecordTripsResponse_RecordFailure
	for i, t := range req.GetTrips() {
		if err := s.recordSingleTrip(ctx, t); err != nil {
			span.RecordError(err, trace.WithAttributes(attribute.Int("trip.idx", i), attribute.String("trip.id", t.GetId())))
			failures = append(failures, &tripspb.RecordTripsResponse_RecordFailure{
				TripId:  t.GetId(),
				Message: err.Error(),
			})
		}
	}

	return &tripspb.RecordTripsResponse{
		Failures: failures,
	}, nil
}

func (s *Server) recordSingleTrip(ctx context.Context, t *tripspb.Trip) error {
	ctx, span := tracer.Start(ctx, "Server.recordSingleTrip",
		trace.WithAttributes(
			attribute.String("trip.id", t.GetId()),
			attribute.String("trip.left_at", t.GetLeftAt()),
			attribute.String("trip.returned_at", t.GetReturnedAt())))
	defer span.End()

	if t.GetId() == "" {
		return status.Errorf(codes.InvalidArgument, "missing ID for trip")
	}
	id, err := uuid.Parse(t.GetId())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid UUID for ID of trip: %s", err)
	}

	if t.GetLeftAt() == "" {
		return status.Error(codes.InvalidArgument, "missing left at time for trip")
	}
	leftAt, err := time.Parse(time.RFC3339, t.GetLeftAt())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid left at time for trip: %s", err)
	}

	if t.GetReturnedAt() == "" {
		return status.Error(codes.InvalidArgument, "missing returned at time for trip")
	}
	returnedAt, err := time.Parse(time.RFC3339, t.GetReturnedAt())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid returned at time for trip: %s", err)
	}

	trip, err := s.q.RecordTrip(ctx, database.RecordTripParams{
		ID:     id,
		LeftAt: leftAt,
		ReturnedAt: sql.NullTime{
			Time:  returnedAt,
			Valid: true,
		},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "recording trip: %s", err)
	}

	if _, err := s.messages.SendTripCompletedMessage(ctx, &messagespb.SendTripCompletedMessageRequest{
		TripId:     trip.ID.String(),
		LeftAt:     trip.LeftAt.Format(time.RFC3339),
		ReturnedAt: trip.ReturnedAt.Time.Format(time.RFC3339),
	}); err != nil {
		return status.Errorf(codes.Internal, "sending trip completed message: %s", err)
	}

	return nil
}
