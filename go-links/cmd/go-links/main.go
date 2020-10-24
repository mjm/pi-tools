package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"go.opentelemetry.io/otel/exporters/metric/prometheus"

	"github.com/mjm/pi-tools/go-links/database"
	_ "github.com/mjm/pi-tools/go-links/service/linksservice"
)

var (
	httpPort = flag.Int("http-port", 4240, "HTTP port to listen on for metrics and API requests")
	dbDSN    = flag.String("db", "dbname=golinks_dev sslmode=disable", "Connection string for connecting to PostgreSQL database for storing links")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	metrics, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		log.Fatalf("Error installing metrics pipeline: %v", err)
	}

	sqlDB, err := sql.Open("postgres", *dbDSN)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	db := database.New(sqlDB)
	if err := db.MigrateIfNeeded(ctx); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	http.Handle("/metrics", metrics)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}
