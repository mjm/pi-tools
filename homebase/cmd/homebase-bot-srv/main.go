package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/homebase/bot/service/messagesservice"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	tripsURL = flag.String("trips-url", "localhost:2121", "URL for trips service to lookup and update trip information")
)

func main() {
	rpc.SetDefaultHTTPPort(6360)
	rpc.SetDefaultGRPCPort(6361)
	flag.Parse()

	stopObs := observability.MustStart("homebase-bot-srv")
	defer stopObs()

	t, err := telegram.New(telegram.Config{
		Token: os.Getenv("TELEGRAM_TOKEN"),
		HTTPClient: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	})
	if err != nil {
		log.Panicf("Error creating Telegram client: %v", err)
	}

	tripsConn := rpc.MustDial(context.Background(), *tripsURL)
	defer tripsConn.Close()

	trips := tripspb.NewTripsServiceClient(tripsConn)

	messagesService := messagesservice.New(t, trips)

	ch := make(chan telegram.UpdateOrError, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.WatchUpdates(ctx, ch, telegram.GetUpdatesRequest{
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
					_, err := t.SendMessage(ctx, telegram.SendMessageRequest{
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
					if err := t.AnswerCallbackQuery(ctx, telegram.AnswerCallbackQueryRequest{
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

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		messagespb.RegisterMessagesServiceServer(s, messagesService)
	}))

	signal.Wait()
}
