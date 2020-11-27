package tripsservice

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/detect-presence/database"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
)

func (s *Server) RecordTrips(ctx context.Context, req *tripspb.RecordTripsRequest) (*tripspb.RecordTripsResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.Int("trip.count", len(req.GetTrips())))

	if len(req.GetTrips()) == 0 {
		return &tripspb.RecordTripsResponse{}, nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "starting transaction: %s", err)
	}
	defer tx.Rollback()

	q := s.q.WithTx(tx)

	var recordedTrips []database.Trip
	for i, t := range req.GetTrips() {
		if t.GetId() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "missing ID for trip %d", i)
		}
		id, err := uuid.Parse(t.GetId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid UUID for ID of trip %d: %s", i, err)
		}

		if t.GetLeftAt() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "missing left at time for trip %d", i)
		}
		leftAt, err := time.Parse(time.RFC3339, t.GetLeftAt())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid left at time for trip %d: %s", i, err)
		}

		if t.GetReturnedAt() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "missing returned at time for trip %d", i)
		}
		returnedAt, err := time.Parse(time.RFC3339, t.GetReturnedAt())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid returned at time for trip %d: %s", i, err)
		}

		trip, err := q.RecordTrip(ctx, database.RecordTripParams{
			ID:     id,
			LeftAt: leftAt,
			ReturnedAt: sql.NullTime{
				Time:  returnedAt,
				Valid: true,
			},
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "recording trip: %s", err)
		}

		recordedTrips = append(recordedTrips, trip)
	}

	if err := tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "recording trips: %s", err)
	}

	for _, trip := range recordedTrips {
		if _, err := s.messages.SendTripCompletedMessage(ctx, &messagespb.SendTripCompletedMessageRequest{
			TripId:     trip.ID.String(),
			LeftAt:     trip.LeftAt.Format(time.RFC3339),
			ReturnedAt: trip.ReturnedAt.Time.Format(time.RFC3339),
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "sending trip completed message: %s", err)
		}
	}

	return &tripspb.RecordTripsResponse{}, nil
}
