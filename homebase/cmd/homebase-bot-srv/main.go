package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/etherlabsio/healthcheck"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/database/migrate"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/homebase/bot/service/messagesservice"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/leader"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/storage"
)

var (
	tripsURL = flag.String("trips-url", "127.0.0.1:2121", "URL for trips service to lookup and update trip information")
	chatID   = flag.Int("chat-id", 223272201, "Chat ID to send messages to")
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/cmd/homebase-bot-srv"

func main() {
	storage.SetDefaultDBName("homebase_bot_dev")
	rpc.SetDefaultHTTPPort(6360)
	rpc.SetDefaultGRPCPort(6361)
	flag.Parse()

	stopObs := observability.MustStart("homebase-bot-srv")
	defer stopObs()

	db := storage.MustOpenDB(migrate.FS)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	election, err := leader.NewElection(leader.Config{
		Key: "service/homebase-bot/leader",
		OnAcquireLeader: func() {
			if err := messagesService.RegisterCommands(ctx); err != nil {
				log.Panicf("Error registering bot commands: %v", err)
			}

			go messagesService.WatchUpdates(ctx)
		},
	})
	if err != nil {
		log.Panicf("Error creating leader election: %v", err)
	}

	metric.Must(otel.Meter(instrumentationName)).NewInt64ValueObserver("homebase.bot.is_leader", func(ctx context.Context, result metric.Int64ObserverResult) {
		var isLeader int64
		if election.IsLeader() {
			isLeader = 1
		}
		result.Observe(isLeader)
	}, metric.WithDescription("Indicates if the instance is the leader and is responsible for watching for incoming messages"))

	go election.Run(ctx)
	defer election.Stop()

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second),
			healthcheck.WithChecker("database", db))))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		messagespb.RegisterMessagesServiceServer(s, messagesService)
	}))

	signal.Wait()
}
