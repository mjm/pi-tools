package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/detect-presence/checker"
	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/detect-presence/service/tripsservice"
	"github.com/mjm/pi-tools/detect-presence/trips"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/rpc"
)

var (
	pingInterval = flag.Duration("ping-interval", 30*time.Second, "How often to check for nearby devices")
	deviceFile   = flag.String("device-file", "", "JSON file to check for device presence instead of using Bluetooth")
	deviceName   = flag.String("device-name", "hci0", "Local Bluetooth device name")
	dbDSN        = flag.String("db", ":memory:", "Connection string for connecting to SQLite3 database for storing trips")
)

var devices = []presence.Device{
	{"callisto", "A0:FB:C5:D3:4D:46"},
	{"matt-watch", "F8:6F:C1:0A:E8:8B"},
}

func main() {
	rpc.SetDefaultHTTPPort(2120)
	flag.Parse()

	stopObs, err := observability.Start("detect-presence-srv")
	if err != nil {
		log.Panicf("Error setting up observability: %v", err)
	}
	defer stopObs()

	db, err := database.Open(*dbDSN)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	if err := db.MigrateIfNeeded(context.Background()); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	tripTracker, err := trips.NewTracker(db)
	if err != nil {
		log.Fatalf("Error setting up trip tracker: %v", err)
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

	go c.Run()

	//http.Handle("/metrics", promhttp.Handler())

	tripsService := tripsservice.New(db)

	log.Fatal(rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		tripspb.RegisterTripsServiceServer(server, tripsService)
	})))
}
