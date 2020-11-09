package messagesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/database"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
)

func (s *Server) SendTripCompletedMessage(ctx context.Context, req *messagespb.SendTripCompletedMessageRequest) (*messagespb.SendTripCompletedMessageResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		label.String("trip.id", req.GetTripId()),
		label.String("trip.left_at", req.GetLeftAt()),
		label.String("trip.returned_at", req.GetReturnedAt()))

	tripID, err := uuid.Parse(req.GetTripId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID for trip ID: %s", err)
	}
	leftAt, err := time.Parse(time.RFC3339, req.GetLeftAt())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid left at timestamp: %s", err)
	}
	returnedAt, err := time.Parse(time.RFC3339, req.GetReturnedAt())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid returned at timestamp: %s", err)
	}
	duration := returnedAt.Sub(leftAt)
	span.SetAttributes(label.Stringer("trip.duration", duration))

	// fetch the most popular three tags for trips and offer them as inline-reply options
	tagsResp, err := s.trips.ListTags(ctx, &tripspb.ListTagsRequest{
		Limit: 3,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fetching popular tags: %v", err)
	}

	var buttonRow []telegram.InlineKeyboardButton
	for _, tag := range tagsResp.GetTags() {
		buttonRow = append(buttonRow, telegram.InlineKeyboardButton{
			Text:         tag.GetName(),
			CallbackData: fmt.Sprintf("TAG_TRIP#%s", tag.GetName()),
		})
	}
	replyMarkup := &telegram.ReplyMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			buttonRow,
		},
	}
	text := fmt.Sprintf("You just returned from a trip that lasted **%s**\\. Do you want to add any tags to the trip?", duration)

	msg, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:      s.chatID,
		Text:        text,
		ParseMode:   telegram.MarkdownV2Mode,
		ReplyMarkup: replyMarkup,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sending message: %s", err)
	}

	span.SetAttributes(label.Int("telegram.message_id", msg.MessageID))

	// before we overwrite it, see if there's an existing message we've already sent about this trip
	existingMessageID, err := s.q.GetMessageForTrip(ctx, tripID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.Internal, "fetching existing message for trip: %s", err)
	}
	span.SetAttributes(label.Int64("telegram.existing_message_id", existingMessageID))

	// record the message ID so we know which trip to update when we get a callback query response
	if err := s.q.SetMessageForTrip(ctx, database.SetMessageForTripParams{
		TripID:    tripID,
		MessageID: int64(msg.MessageID),
	}); err != nil {
		return nil, status.Errorf(codes.Internal, "saving message ID for trip: %s", err)
	}

	// if there was an existing message from when the trip began, delete that message from the chat.
	// we do this instead of updating the existing message because we want there to be a notification
	// for this message.
	if existingMessageID != 0 {
		if err := s.t.DeleteMessage(ctx, telegram.DeleteMessageRequest{
			ChatID:    s.chatID,
			MessageID: int(existingMessageID),
		}); err != nil {
			return nil, status.Errorf(codes.Internal, "deleting previous message: %s", err)
		}
	}

	return &messagespb.SendTripCompletedMessageResponse{}, nil
}
