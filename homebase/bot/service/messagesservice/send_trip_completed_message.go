package messagesservice

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
)

func (s *Server) SendTripCompletedMessage(ctx context.Context, req *messagespb.SendTripCompletedMessageRequest) (*messagespb.SendTripCompletedMessageResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		label.String("trip.id", req.GetTripId()),
		label.String("trip.left_at", req.GetLeftAt()),
		label.String("trip.returned_at", req.GetReturnedAt()))

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
			CallbackData: tag.GetName(),
		})
	}

	text := fmt.Sprintf("You just returned from a trip that lasted **%s**\\. Do you want to add any tags to the trip?", duration)
	_, err = s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:    s.chatID,
		Text:      text,
		ParseMode: telegram.MarkdownV2Mode,
		ReplyMarkup: &telegram.ReplyMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				buttonRow,
			},
		},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sending message: %s", err)
	}
	return &messagespb.SendTripCompletedMessageResponse{}, nil
}
