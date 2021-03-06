package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/etherlabsio/healthcheck"
	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/detect-presence/checker"
	"github.com/mjm/pi-tools/detect-presence/database/migrate"
	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/detect-presence/service/appservice"
	"github.com/mjm/pi-tools/detect-presence/service/tripsservice"
	"github.com/mjm/pi-tools/detect-presence/trips"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/storage"
)

var (
	pingInterval    = flag.Duration("ping-interval", 30*time.Second, "How often to check for nearby devices")
	deviceFile      = flag.String("device-file", "", "JSON file to check for device presence instead of using Bluetooth")
	deviceName      = flag.String("device-name", "hci0", "Local Bluetooth device name")
	messagesURL     = flag.String("messages-url", "127.0.0.1:6361", "URL for messages service to use to send chat messages")
	mode            = flag.String("mode", "server", "Mode (server or client) to use to detect presence")
	githubTokenPath = flag.String("github-token-path", "/secrets/github-token", "Path to file containing GitHub PAT token")
)

var devices = []presence.Device{
	{Name: "canary", Addr: "20:7D:74:20:C7:FD", Canary: true},
	{Name: "callisto", Addr: "A0:FB:C5:D3:4D:46"},
	{Name: "matt-watch", Addr: "F8:6F:C1:0A:E8:8B"},
}

func main() {
	rpc.SetDefaultHTTPPort(2120)
	rpc.SetDefaultGRPCPort(2121)
	storage.SetDefaultDBName("presence_dev")
	flag.Parse()

	if *mode != "server" && *mode != "client" {
		log.Panicf("invalid mode %q", *mode)
	}

	stopObs := observability.MustStart("detect-presence-srv")
	defer stopObs()

	db := storage.MustOpenDB(migrate.FS)

	messagesConn := rpc.MustDial(context.Background(), *messagesURL)
	defer messagesConn.Close()

	messages := messagespb.NewMessagesServiceClient(messagesConn)

	if *mode == "server" {
		tripTracker, err := trips.NewTracker(db, messages)
		if err != nil {
			log.Panicf("Error setting up trip tracker: %v", err)
		}

		t := presence.NewTracker()
		t.OnLeave(tripTracker)
		t.OnReturn(tripTracker)

		var d detector.Detector
		if *deviceFile != "" {
			d = &detector.FileDetector{
				Path: *deviceFile,
			}
		} else {
			d = &detector.HCIDetector{
				DeviceName: *deviceName,
			}
		}

		c := &checker.Checker{
			Tracker:  t,
			Detector: d,
			Interval: *pingInterval,
			Devices:  devices,
		}
		go c.Run(context.Background(), nil)
	}

	tokenData, err := ioutil.ReadFile(*githubTokenPath)
	if err != nil {
		log.Panicf("reading github token: %v", err)
	}

	token := &oauth2.Token{AccessToken: strings.TrimSpace(string(tokenData))}
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)
	githubClient := github.NewClient(httpClient)

	appService := appservice.New(githubClient)
	http.Handle("/app/download",
		otelhttp.WithRouteTag("DownloadApp", http.HandlerFunc(appService.DownloadApp)))
	http.Handle("/app/install",
		otelhttp.WithRouteTag("InstallApp", http.HandlerFunc(appService.InstallApp)))
	http.Handle("/app/install_manifest",
		otelhttp.WithRouteTag("InstallManifest", http.HandlerFunc(appService.InstallManifest)))
	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second),
			healthcheck.WithChecker("database", db))))

	tripsService := tripsservice.New(db, messages)
	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		tripspb.RegisterTripsServiceServer(server, tripsService)
	}))

	signal.Wait()
}
