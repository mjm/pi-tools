package messagesservice

import (
	"context"
	"log"

	"github.com/mjm/pi-tools/homebase/bot/telegram"
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
	log.Printf("Got callback data %s", cbq.Data)
	if err := s.t.AnswerCallbackQuery(ctx, telegram.AnswerCallbackQueryRequest{
		CallbackQueryID: cbq.ID,
		Text:            "Got the message!",
	}); err != nil {
		return err
	}
	return nil
}
