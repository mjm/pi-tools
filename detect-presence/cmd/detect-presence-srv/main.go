package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mjm/pi-tools/detect-presence/checker"
	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/presence"
	"github.com/mjm/pi-tools/detect-presence/trips"
)

var (
	httpPort     = flag.Int("http-port", 2120, "HTTP port to listen on for metrics and API requests")
	pingInterval = flag.Duration("ping-interval", 30*time.Second, "How often to check for nearby devices")
	deviceName   = flag.String("device-name", "hci0", "Local Bluetooth device name")
	dbDSN        = flag.String("db", ":memory:", "Connection string for connecting to SQLite3 database for storing trips")
)

var devices = []presence.Device{
	{"titan", "C4:98:80:8C:6C:2D"},
	{"matt-watch", "F8:6F:C1:0A:E8:8B"},
}

func main() {
	flag.Parse()

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

	c := &checker.Checker{
		Tracker:    t,
		Interval:   *pingInterval,
		DeviceName: *deviceName,
		Devices:    devices,
	}

	go c.Run()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}
