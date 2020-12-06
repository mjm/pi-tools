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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

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
	tripsURL   = flag.String("trips-url", "localhost:2121", "URL for trips service to lookup and update trip information")
	chatID     = flag.Int("chat-id", 223272201, "Chat ID to send messages to")
	namespace  = flag.String("namespace", "", "Kubernetes namespace to use for leader election")
	instanceID = flag.String("instance-id", "", "Name of this instance of the service to use for leader election")
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/cmd/homebase-bot-srv"

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

	if *instanceID == "" {
		isLeader = 1
		if err := messagesService.RegisterCommands(ctx); err != nil {
			log.Panicf("Error registering bot commands: %v", err)
		}

		go messagesService.WatchUpdates(ctx)
	} else {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			log.Panicf("Error creating in-cluster Kubernetes config: %v", err)
		}
		clientset := kubernetes.NewForConfigOrDie(cfg)
		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Namespace: *namespace,
				Name:      "homebase-bot-srv-leader",
			},
			Client: clientset.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: *instanceID,
			},
		}

		go leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
			Lock:            lock,
			ReleaseOnCancel: true,
			LeaseDuration:   15 * time.Second,
			RenewDeadline:   10 * time.Second,
			RetryPeriod:     2 * time.Second,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					isLeader = 1

					log.Printf("Started leading")
					if err := messagesService.RegisterCommands(ctx); err != nil {
						log.Printf("Error registering bot commands: %v", err)
					}

					messagesService.WatchUpdates(ctx)
				},
				OnStoppedLeading: func() {
					log.Fatalf("Stopped leading")
				},
				OnNewLeader: func(identity string) {
					if identity == *instanceID {
						return
					}
					log.Printf("New leader: %s", identity)
				},
			},
		})
	}

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second),
			healthcheck.WithChecker("database", db))))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		messagespb.RegisterMessagesServiceServer(s, messagesService)
	}))

	signal.Wait()
}
