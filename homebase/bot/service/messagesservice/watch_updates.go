package messagesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) WatchUpdates(ctx context.Context) {
	ch := make(chan telegram.UpdateOrError, 10)
	s.t.WatchUpdates(ctx, ch, telegram.GetUpdatesRequest{
		Timeout: 30,
	})

	for updateOrErr := range ch {
		if updateOrErr.Err != nil {
			s.metrics.TelegramUpdateErrorsTotal.Add(ctx, 1)
			log.Printf("Error getting updates: %v", updateOrErr.Err)
		} else {
			update := updateOrErr.Update
			if update.CallbackQuery != nil {
				s.metrics.TelegramUpdateTotal.Add(ctx, 1,
					label.String("update_type", "callback_query"))

				if err := s.handleCallbackQuery(ctx, update.CallbackQuery); err != nil {
					log.Printf("Error answering callback query: %v", err)
				}
			} else if update.Message != nil {
				s.metrics.TelegramUpdateTotal.Add(ctx, 1,
					label.String("update_type", "message"))

				if err := s.handleMessage(ctx, update.Message); err != nil {
					log.Printf("Error responding to message: %v", err)
				}
			} else {
				log.Printf("Received update: %+v", updateOrErr.Update)
			}
		}
	}
}

func (s *Server) handleMessage(ctx context.Context, msg *telegram.Message) error {
	ctx, span := tracer.Start(ctx, "MessagesService.handleMessage",
		trace.WithAttributes(
			label.Int("telegram.message_id", msg.MessageID)))
	defer span.End()

	// check if the message is a command
	if strings.HasPrefix(msg.Text, "/") {
		if strings.HasPrefix(msg.Text, "/ignore") {
			if err := s.handleIgnoreCommand(ctx, msg); err != nil {
				// TODO respond to the user with the error
				return err
			}
		} else if strings.HasPrefix(msg.Text, "/tag") {
			if err := s.handleTagCommand(ctx, msg); err != nil {
				// TODO respond to the user with the error
				return err
			}
		} else if strings.HasPrefix(msg.Text, "/untag") {
			if err := s.handleUntagCommand(ctx, msg); err != nil {
				// TODO respond to the user with the error
				return err
			}
		}
	}

	return nil
}

func (s *Server) handleCallbackQuery(ctx context.Context, cbq *telegram.CallbackQuery) error {
	ctx, span := tracer.Start(ctx, "MessagesService.handleCallbackQuery",
		trace.WithAttributes(
			label.Int("telegram.message_id", cbq.Message.MessageID),
			label.String("telegram.callback_query.data", cbq.Data)))
	defer span.End()

	if strings.HasPrefix(cbq.Data, "TAG_TRIP#") {
		tagName := cbq.Data[9:]
		span.SetAttributes(label.String("tag.name", tagName))

		tripID, err := s.q.GetTripForMessage(ctx, int64(cbq.Message.MessageID))
		if err != nil {
			span.RecordError(ctx, err)
			var text string
			if errors.Is(err, sql.ErrNoRows) {
				text = "Sorry, I couldn't find the trip that goes with that message."
			} else {
				text = fmt.Sprintf("Sorry, something unexpected happened: %s", err)
			}

			if err := s.t.AnswerCallbackQuery(ctx, telegram.AnswerCallbackQueryRequest{
				CallbackQueryID: cbq.ID,
				Text:            text,
			}); err != nil {
				return err
			}

			return err
		}

		_, err = s.trips.UpdateTripTags(ctx, &tripspb.UpdateTripTagsRequest{
			TripId:    tripID.String(),
			TagsToAdd: []string{tagName},
		})
		if err != nil {
			// TODO respond to user
			span.RecordError(ctx, err)
			return err
		}

		if err := s.t.AnswerCallbackQuery(ctx, telegram.AnswerCallbackQueryRequest{
			CallbackQueryID: cbq.ID,
			Text:            fmt.Sprintf("Done! I added the %q tag to that trip.", tagName),
		}); err != nil {
			return spanerr.RecordError(ctx, err)
		}
		return nil
	}

	if err := s.t.AnswerCallbackQuery(ctx, telegram.AnswerCallbackQueryRequest{
		CallbackQueryID: cbq.ID,
		Text:            "Sorry, I don't know what to do with this.",
	}); err != nil {
		return spanerr.RecordError(ctx, err)
	}
	return nil
}

func (s *Server) handleIgnoreCommand(ctx context.Context, msg *telegram.Message) error {
	ctx, span := tracer.Start(ctx, "MessagesService.handleIgnoreCommand",
		trace.WithAttributes(label.String("telegram.message.text", msg.Text)))
	defer span.End()

	resp, err := s.trips.GetLastCompletedTrip(ctx, &tripspb.GetLastCompletedTripRequest{})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	_, err = s.trips.IgnoreTrip(ctx, &tripspb.IgnoreTripRequest{
		Id: resp.GetTrip().GetId(),
	})

	returnedAt, err := time.Parse(time.RFC3339, resp.GetTrip().GetReturnedAt())
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             fmt.Sprintf("Done! Your trip from %s ago has been ignored.", time.Now().Sub(returnedAt)),
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

	resp, err := s.trips.GetLastCompletedTrip(ctx, &tripspb.GetLastCompletedTripRequest{})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.trips.UpdateTripTags(ctx, &tripspb.UpdateTripTagsRequest{
		TripId:    resp.GetTrip().GetId(),
		TagsToAdd: parseTagList(strings.TrimPrefix(msg.Text, "/tag")),
	}); err != nil {
		return err
	}

	returnedAt, err := time.Parse(time.RFC3339, resp.GetTrip().GetReturnedAt())
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             fmt.Sprintf("Done! Your trip from %s ago has been tagged.", time.Now().Sub(returnedAt)),
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

	resp, err := s.trips.GetLastCompletedTrip(ctx, &tripspb.GetLastCompletedTripRequest{})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.trips.UpdateTripTags(ctx, &tripspb.UpdateTripTagsRequest{
		TripId:       resp.GetTrip().GetId(),
		TagsToRemove: parseTagList(strings.TrimPrefix(msg.Text, "/untag")),
	}); err != nil {
		return err
	}

	returnedAt, err := time.Parse(time.RFC3339, resp.GetTrip().GetReturnedAt())
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	if _, err := s.t.SendMessage(ctx, telegram.SendMessageRequest{
		ChatID:           msg.Chat.ID,
		Text:             fmt.Sprintf("Done! Your trip from %s ago has been untagged.", time.Now().Sub(returnedAt)),
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
