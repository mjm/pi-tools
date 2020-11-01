package messagesservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

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
			log.Printf("Error getting updates: %v", updateOrErr.Err)
		} else {
			update := updateOrErr.Update
			if update.CallbackQuery != nil {
				if err := s.handleCallbackQuery(ctx, update.CallbackQuery); err != nil {
					log.Printf("Error answering callback query: %v", err)
				}
			} else {
				log.Printf("Received update: %+v", updateOrErr.Update)
			}
		}
	}
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
