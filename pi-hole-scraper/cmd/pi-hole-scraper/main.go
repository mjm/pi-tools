package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/mjm/pi-tools/pi-hole-scraper/scraper"
)

var (
	piHoleAddress  = flag.String("pi-hole-address", "http://localhost", "Address of the PiHole server to scrape from")
	piHolePassword = flag.String("pi-hole-password", "", "Password for accessing PiHole API (WEBPASSWORD in setupVars.conf)")

	influxDBAddress  = flag.String("influxdb-address", "http://localhost:8086", "Address of the InfluxDB instance to send the metrics to")
	influxDBToken    = flag.String("influxdb-token", "", "Authentication token for accessing InfluxDB")
	influxDBDatabase = flag.String("influxdb-database", "pihole", "Database name to send measurements to")
)

func main() {
	flag.Parse()

	s := scraper.New(*piHoleAddress, *piHolePassword)
	ctx := context.Background()

	stats, err := s.GetStats(ctx)
	if err != nil {
		log.Fatalf("Error scraping stats from PiHole: %v", err)
	}

	log.Printf("Stats: %#v", stats)

	influx := influxdb2.NewClient(*influxDBAddress, *influxDBToken)
	influxWrite := influx.WriteAPIBlocking("", *influxDBDatabase)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Error getting hostname: %v", err)
	}
	point := write.NewPoint("pihole.stats", map[string]string{
		"host": hostname,
	}, map[string]interface{}{
		"domains_blocked":      stats.DomainsBeingBlocked,
		"dns_queries_today":    stats.DNSQueriesToday,
		"ads_percentage_today": stats.AdsPercentageToday,
		"ads_blocked_today":    stats.AdsBlockedToday,
	}, time.Now())

	if err := influxWrite.WritePoint(ctx, point); err != nil {
		log.Fatalf("Error writing stats to InfluxDB: %v", err)
	}

	log.Print("Successfully sent PiHole stats to InfluxDB")
}
