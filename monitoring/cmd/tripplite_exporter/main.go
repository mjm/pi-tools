package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/etherlabsio/healthcheck"
	"github.com/zserge/hid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

func main() {
	rpc.SetDefaultHTTPPort(8080)
	flag.Parse()
	hid.Logger = log.New(os.Stderr, "hid", log.LstdFlags)

	stopObs := observability.MustStart("tripplite-exporter")
	defer stopObs()

	var device hid.Device
	hid.UsbWalk(func(d hid.Device) {
		info := d.Info()
		log.Printf("Found device %04x:%04x", info.Vendor, info.Product)
		if info.Vendor == 0x09ae {
			log.Printf("This device is a TrippLite device, using it")
			device = d
		}
	})

	if device == nil {
		log.Panicf("No TrippLite devices found")
	}

	if err := device.Open(); err != nil {
		log.Panicf("opening device: %v", err)
	}
	defer device.Close()

	ObserveDevices(context.Background(), []hid.Device{device})

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))
	rpc.ListenAndServe()

	signal.Wait()
}
