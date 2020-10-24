package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/exporters/metric/prometheus"

	_ "github.com/mjm/pi-tools/go-links/service/linksservice"
)

var (
	httpPort = flag.Int("http-port", 4240, "HTTP port to listen on for metrics and API requests")
	dbDSN    = flag.String("db", ":memory:", "Connection string for connecting to SQLite3 database for storing links")
)

func main() {
	flag.Parse()

	metrics, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Fatalf("Error installing metrics pipeline: %v", err)
	}

	//db, err := database.Open(*dbDSN)
	//if err != nil {
	//	log.Fatalf("Error opening database: %v", err)
	//}
	//if err := db.MigrateIfNeeded(context.Background()); err != nil {
	//	log.Fatalf("Error migrating database: %v", err)
	//}

	http.Handle("/metrics", metrics)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}
