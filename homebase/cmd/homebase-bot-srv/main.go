package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/etherlabsio/healthcheck"
	"github.com/hashicorp/consul/api"
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
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/storage"
)

var (
	tripsURL    = flag.String("trips-url", "localhost:2121", "URL for trips service to lookup and update trip information")
	chatID      = flag.Int("chat-id", 223272201, "Chat ID to send messages to")
	leaderElect = flag.Bool("leader-elect", false, "Enable leader election using Consul")
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/cmd/homebase-bot-srv"

const leaderKeyName = "service/homebase-bot/leader"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var isLeader int64
	metric.Must(otel.Meter(instrumentationName)).NewInt64ValueObserver("homebase.bot.is_leader", func(ctx context.Context, result metric.Int64ObserverResult) {
		result.Observe(isLeader)
	}, metric.WithDescription("Indicates if the instance is the leader and is responsible for watching for incoming messages"))

	stopLeader := func() {}
	if !*leaderElect {
		log.Printf("No leader election desired.")

		isLeader = 1
		if err := messagesService.RegisterCommands(ctx); err != nil {
			log.Panicf("Error registering bot commands: %v", err)
		}

		go messagesService.WatchUpdates(ctx)
	} else {
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			log.Panicf("Error creating Consul client: %v", err)
		}

		log.Printf("Creating Consul session for key %s", leaderKeyName)
		sessionID, _, err := client.Session().Create(&api.SessionEntry{
			Name: leaderKeyName,
			TTL:  "15s",
		}, nil)
		if err != nil {
			log.Panicf("Error creating Consul session: %v", err)
		}
		go client.Session().RenewPeriodic("15s", sessionID, nil, ctx.Done())

		go func() {
			gotLeader, _, err := client.KV().Acquire(&api.KVPair{
				Key:     leaderKeyName,
				Value:   []byte(sessionID),
				Session: sessionID,
			}, nil)
			if err != nil {
				log.Panicf("Error trying to acquire leadership: %v", err)
			}

			var waitIndex uint64
			for !gotLeader {
				// Watch for the key to change to not have a session anymore, and try to grab leader
				kvPair, meta, err := client.KV().Get(leaderKeyName, &api.QueryOptions{
					WaitIndex: waitIndex,
				})
				if err != nil {
					log.Printf("Error checking on leader key %q: %v", leaderKeyName, err)
					waitIndex = 0
					time.Sleep(10 * time.Second)
					continue
				}

				if meta.LastIndex < waitIndex {
					waitIndex = 0
				} else {
					waitIndex = meta.LastIndex
				}

				if kvPair == nil || kvPair.Session == "" {
					// There's no longer a session, so try to acquire leadership
					gotLeader, _, err = client.KV().Acquire(&api.KVPair{
						Key:     leaderKeyName,
						Value:   []byte(sessionID),
						Session: sessionID,
					}, nil)
					if err != nil {
						log.Panicf("Error trying to acquire leadership: %v", err)
					}
				}
			}

			// Once we get here, we are the leader
			isLeader = 1

			go func() {
				log.Printf("Started leading")
				if err := messagesService.RegisterCommands(ctx); err != nil {
					log.Printf("Error registering bot commands: %v", err)
				}

				messagesService.WatchUpdates(ctx)
			}()

			stopLeader = func() {
				log.Printf("Releasing leadership voluntarily")
				_, _, err = client.KV().Release(&api.KVPair{
					Key:     leaderKeyName,
					Value:   []byte{},
					Session: sessionID,
				}, nil)
				if err != nil {
					log.Printf("Failed to release lock: %v", err)
				}
			}

			for {
				// Watch to see if our leadership gets revoked
				kvPair, meta, err := client.KV().Get(leaderKeyName, &api.QueryOptions{
					WaitIndex: waitIndex,
				})
				if err != nil {
					log.Panicf("Error checking on leader key %q: %v", leaderKeyName, err)
				}

				if meta.LastIndex < waitIndex {
					waitIndex = 0
				} else {
					waitIndex = meta.LastIndex
				}

				if kvPair.Session != sessionID {
					log.Fatalf("Lost leadership")
				}
			}
		}()
	}

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second),
			healthcheck.WithChecker("database", db))))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		messagespb.RegisterMessagesServiceServer(s, messagesService)
	}))

	signal.Wait()
	stopLeader()
}
