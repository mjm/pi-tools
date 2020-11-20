package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	if *instanceID == "" {
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

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		messagespb.RegisterMessagesServiceServer(s, messagesService)
	}))

	signal.Wait()
}
