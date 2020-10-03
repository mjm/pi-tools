package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mjm/pi-tools/detect-presence/checker"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

var (
	httpPort     = flag.Int("http-port", 2120, "HTTP port to listen on for metrics and API requests")
	pingInterval = flag.Duration("ping-interval", 30*time.Second, "How often to check for nearby devices")
	deviceName   = flag.String("device-name", "hci0", "Local Bluetooth device name")
)

var devices = []presence.Device{
	{"titan", "C4:98:80:8C:6C:2D"},
	{"matt-watch", "F8:6F:C1:0A:E8:8B"},
}

func main() {
	flag.Parse()

	t := presence.NewTracker()
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
