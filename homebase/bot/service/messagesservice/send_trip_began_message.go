package messagesservice

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/homebase/bot/database"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
)

func (s *Server) SendTripBeganMessage(ctx context.Context, req *messagespb.SendTripBeganMessageRequest) (*messagespb.SendTripBeganMessageResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		label.String("trip.id", req.GetTripId()),
		label.String("trip.left_at", req.GetLeftAt()))

	tripID, err := uuid.Parse(req.GetTripId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID for trip ID: %s", err)
	}
	//leftAt, err := time.Parse(time.RFC3339, req.GetLeftAt())
	//if err != nil {
	//	return nil, status.Errorf(codes.InvalidArgument, "invalid left at timestamp: %s", err)
	//}

	msg, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID: s.chatID,
		Text:   "It looks like you've left home.",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sending message: %s", err)
	}

	span.SetAttributes(label.Int("telegram.message_id", msg.MessageID))

	// record the message ID so we know which trip to update when we get a callback query response
	if err := s.q.SetMessageForTrip(ctx, database.SetMessageForTripParams{
		TripID:    tripID,
		MessageID: int64(msg.MessageID),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "saving message ID for trip: %s", err)
	}

	return &messagespb.SendTripBeganMessageResponse{}, nil
}
