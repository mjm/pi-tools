package messagesservice

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/pkg/spanerr"
)

var commands = []telegram.BotCommand{
	{
		Command:     "tag",
		Description: "Add one or more tags to the most recent trip.",
	},
	{
		Command:     "untag",
		Description: "Remove one or more tags from the most recent trip.",
	},
	{
		Command:     "ignore",
		Description: "Ignore the most recent trip, hiding it from the list of trips in Homebase.",
	},
}

func (s *Server) RegisterCommands(ctx context.Context) error {
	if err := s.t.SetMyCommands(ctx, telegram.SetMyCommandsRequest{
		Commands: commands,
	}); err != nil {
		return status.Errorf(codes.Internal, "setting commands: %s", err)
	}

	return nil
}

func (s *Server) handleIgnoreCommand(ctx context.Context, msg *telegram.Message) error {
	ctx, span := tracer.Start(ctx, "MessagesService.handleIgnoreCommand",
		trace.WithAttributes(label.String("telegram.message.text", msg.Text)))
	defer span.End()

	resp, err := s.trips.GetLastCompletedTrip(ctx, &trips.GetLastCompletedTripRequest{})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	_, err = s.trips.IgnoreTrip(ctx, &trips.IgnoreTripRequest{
		Id: resp.GetTrip().GetId(),
	})

	returnedAt, err := time.Parse(time.RFC3339, resp.GetTrip().GetReturnedAt())
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	var text strings.Builder
	if err := templates.ExecuteTemplate(&text, tripIgnoredTemplate, &tripIgnoredTemplateInput{
		ReturnedAt: returnedAt,
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             text.String(),
		ReplyToMessageID: msg.MessageID,
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	return nil
}

func (s *Server) handleTagCommand(ctx context.Context, msg *telegram.Message) error {
	ctx, span := tracer.Start(ctx, "MessagesService.handleTagCommand",
		trace.WithAttributes(label.String("telegram.message.text", msg.Text)))
	defer span.End()

	resp, err := s.trips.GetLastCompletedTrip(ctx, &trips.GetLastCompletedTripRequest{})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.trips.UpdateTripTags(ctx, &trips.UpdateTripTagsRequest{
		TripId:    resp.GetTrip().GetId(),
		TagsToAdd: parseTagList(strings.TrimPrefix(msg.Text, "/tag")),
	}); err != nil {
		return err
	}

	returnedAt, err := time.Parse(time.RFC3339, resp.GetTrip().GetReturnedAt())
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	var text strings.Builder
	if err := templates.ExecuteTemplate(&text, tripTaggedTemplate, &tripTaggedTemplateInput{
		ReturnedAt: returnedAt,
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             text.String(),
		ReplyToMessageID: msg.MessageID,
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	return nil
}

func (s *Server) handleUntagCommand(ctx context.Context, msg *telegram.Message) error {
	ctx, span := tracer.Start(ctx, "MessagesService.handleUntagCommand",
		trace.WithAttributes(label.String("telegram.message.text", msg.Text)))
	defer span.End()

	resp, err := s.trips.GetLastCompletedTrip(ctx, &trips.GetLastCompletedTripRequest{})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.trips.UpdateTripTags(ctx, &trips.UpdateTripTagsRequest{
		TripId:       resp.GetTrip().GetId(),
		TagsToRemove: parseTagList(strings.TrimPrefix(msg.Text, "/untag")),
	}); err != nil {
		return err
	}

	returnedAt, err := time.Parse(time.RFC3339, resp.GetTrip().GetReturnedAt())
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	var text strings.Builder
	if err := templates.ExecuteTemplate(&text, tripUntaggedTemplate, &tripUntaggedTemplateInput{
		ReturnedAt: returnedAt,
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             text.String(),
		ReplyToMessageID: msg.MessageID,
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	return nil
}

func parseTagList(s string) []string {
	tagList := strings.Split(strings.TrimSpace(s), ",")

	var tags []string
	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
