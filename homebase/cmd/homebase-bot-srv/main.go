package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/database/migrate"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/homebase/bot/service/messagesservice"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/storage"
)

var (
	tripsURL = flag.String("trips-url", "localhost:2121", "URL for trips service to lookup and update trip information")
	chatID   = flag.Int("chat-id", 223272201, "Chat ID to send messages to")
)

func main() {
	storage.SetDefaultDBName("homebase_bot_dev")
	rpc.SetDefaultHTTPPort(6360)
	rpc.SetDefaultGRPCPort(6361)
	flag.Parse()

	stopObs := observability.MustStart("homebase-bot-srv")
	defer stopObs()

	db := storage.MustOpenDB(migrate.Data)

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

	messagesService := messagesservice.New(db, t, trips, *chatID)

	watchCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go messagesService.WatchUpdates(watchCtx)

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		messagespb.RegisterMessagesServiceServer(s, messagesService)
	}))

	signal.Wait()
}
