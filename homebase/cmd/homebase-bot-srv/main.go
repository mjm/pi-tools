package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
)

func main() {
	flag.Parse()

	stopObs, err := observability.Start("homebase-bot-srv")
	if err != nil {
		log.Panicf("Error setting up observability: %v", err)
	}
	defer stopObs()

	c, err := telegram.New(telegram.Config{
		Token: os.Getenv("TELEGRAM_TOKEN"),
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	})
	if err != nil {
		log.Panicf("Error creating Telegram client: %v", err)
	}

	ch := make(chan telegram.UpdateOrError, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.WatchUpdates(ctx, ch, telegram.GetUpdatesRequest{
		Timeout: 30,
	})

	go func() {
		for updateOrErr := range ch {
			if updateOrErr.Err != nil {
				log.Printf("Error getting updates: %v", updateOrErr.Err)
			} else {
				update := updateOrErr.Update
				if update.Message != nil {
					log.Printf("Received message: %#v", update.Message)
					_, err := c.SendMessage(ctx, telegram.SendMessageRequest{
						ChatID:           update.Message.Chat.ID,
						Text:             fmt.Sprintf("Received message: %s", update.Message.Text),
						ReplyToMessageID: update.Message.MessageID,
						ReplyMarkup: &telegram.ReplyMarkup{
							InlineKeyboard: [][]telegram.InlineKeyboardButton{
								{
									{
										Text:         "dog walk",
										CallbackData: "dog walk",
									},
									{
										Text:         "package pickup",
										CallbackData: "package pickup",
									},
								},
							},
						},
					})
					if err != nil {
						log.Printf("Error sending reply: %v", err)
					}
				} else if update.CallbackQuery != nil {
					log.Printf("Got callback data %s", update.CallbackQuery.Data)
					if err := c.AnswerCallbackQuery(ctx, telegram.AnswerCallbackQueryRequest{
						CallbackQueryID: update.CallbackQuery.ID,
						Text:            "Got the message!",
					}); err != nil {
						log.Printf("Error answering callback query: %v", err)
					}
				} else {
					log.Printf("Received update: %+v", updateOrErr.Update)
				}
			}
		}
	}()

	signal.Wait()
}
