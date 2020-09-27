package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/mjm/pi-tools/detect-presence/detector"
)

var (
	influxDBAddress  = flag.String("influxdb-address", "http://localhost:8086", "Address of the InfluxDB instance to send the metrics to")
	influxDBToken    = flag.String("influxdb-token", "", "Authentication token for accessing InfluxDB")
	influxDBDatabase = flag.String("influxdb-database", "pihole", "Database name to send measurements to")
)

var devices = map[string]string{
	"titan":      "C4:98:80:8C:6C:2D",
	"matt-watch": "F8:6F:C1:0A:E8:8B",
}

var influxClient api.WriteAPIBlocking
var hostname string

func main() {
	flag.Parse()
	ctx := context.Background()

	influx := influxdb2.NewClient(*influxDBAddress, *influxDBToken)
	influxClient = influx.WriteAPIBlocking("", *influxDBDatabase)

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %v", err)
	}

	for name, addr := range devices {
		detectAndReport(ctx, name, addr)
	}
}

func detectAndReport(ctx context.Context, name, addr string) {
	log.Printf("Detecting device %q (%s)", name, addr)

	detectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	present, err := detector.DetectDevice(detectCtx, addr)
	if err != nil {
		log.Printf("Error detecting device %q: %v", name, err)
		point := write.NewPoint("device_presence.error", map[string]string{
			"host":        hostname,
			"device_name": name,
			"device_addr": addr,
		}, map[string]interface{}{
			"error": 1,
		}, time.Now())
		if err := influxClient.WritePoint(ctx, point); err != nil {
			log.Printf("Error writing to InfluxDB: %v", err)
		}
		return
	}

	// coerce to an integer because grafana doesn't like influx booleans
	var presentInt int
	if present {
		presentInt = 1
	}

	point := write.NewPoint("device_presence", map[string]string{
		"host":        hostname,
		"device_name": name,
		"device_addr": addr,
	}, map[string]interface{}{
		"present": presentInt,
	}, time.Now())

	if err := influxClient.WritePoint(ctx, point); err != nil {
		log.Printf("Error writing to InfluxDB: %v", err)
		return
	}

	log.Printf("Successfully reported presence for %s", name)
}
