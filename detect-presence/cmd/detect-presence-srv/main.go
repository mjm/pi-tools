package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/detect-presence/checker"
	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/detect-presence/service/tripsservice"
	"github.com/mjm/pi-tools/detect-presence/trips"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
	"github.com/mjm/pi-tools/rpc"
)

var (
	pingInterval = flag.Duration("ping-interval", 30*time.Second, "How often to check for nearby devices")
	deviceFile   = flag.String("device-file", "", "JSON file to check for device presence instead of using Bluetooth")
	deviceName   = flag.String("device-name", "hci0", "Local Bluetooth device name")
	dbDSN        = flag.String("db", "dbname=presence_dev sslmode=disable", "Connection string for connecting to PostgreSQL database for storing trips")
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

	sqlDB, err := sql.Open("postgres", *dbDSN)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	db := otelsql.NewDBWithTracing(sqlDB,
		otelsql.WithAttributes(
			semconv.DBSystemPostgres,
			// assuming this is safe to include since it was on the command-line.
			// passwords should come from a file or environment variable.
			semconv.DBConnectionStringKey.String(*dbDSN)))

	if err := database.New(db).MigrateIfNeeded(context.Background()); err != nil {
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

	tripsService := tripsservice.New(db)

	log.Fatal(rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		tripspb.RegisterTripsServiceServer(server, tripsService)
	})))
}
