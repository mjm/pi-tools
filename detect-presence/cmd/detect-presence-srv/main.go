package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/detect-presence/checker"
	"github.com/mjm/pi-tools/detect-presence/database/migrate"
	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/detect-presence/service/tripsservice"
	"github.com/mjm/pi-tools/detect-presence/trips"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/storage"
)

var (
	pingInterval = flag.Duration("ping-interval", 30*time.Second, "How often to check for nearby devices")
	deviceFile   = flag.String("device-file", "", "JSON file to check for device presence instead of using Bluetooth")
	deviceName   = flag.String("device-name", "hci0", "Local Bluetooth device name")
	messagesURL  = flag.String("messages-url", "localhost:6361", "URL for messages service to use to send chat messages")
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

	stopObs := observability.MustStart("detect-presence-srv")
	defer stopObs()

	db := storage.MustOpenDB(migrate.Data)

	messagesConn := rpc.MustDial(context.Background(), *messagesURL)
	defer messagesConn.Close()

	messages := messagespb.NewMessagesServiceClient(messagesConn)

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

	tripsService := tripsservice.New(db)
	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		tripspb.RegisterTripsServiceServer(server, tripsService)
	}))

	signal.Wait()
}
