package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
	"github.com/mjm/pi-tools/go-links/service/linksservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
	"github.com/mjm/pi-tools/rpc"
)

var (
	dbDSN = flag.String("db", "dbname=golinks_dev sslmode=disable", "Connection string for connecting to PostgreSQL database for storing links")
)

func main() {
	rpc.SetDefaultHTTPPort(4240)
	flag.Parse()

	ctx := context.Background()

	stopObs, err := observability.Start()
	if err != nil {
		log.Panicf("Error setting up observability: %v", err)
	}
	defer stopObs()

	sqlDB, err := sql.Open("postgres", *dbDSN)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	db := database.New(otelsql.NewDBWithTracing(sqlDB,
		otelsql.WithAttributes(
			semconv.DBSystemPostgres,
			// assuming this is safe to include since it was on the command-line.
			// passwords should come from a file or environment variable.
			semconv.DBConnectionStringKey.String(*dbDSN))))

	if err := db.MigrateIfNeeded(ctx); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	linksService := linksservice.New(db)
	http.Handle("/", otelhttp.WithRouteTag("HandleShortLink", http.HandlerFunc(linksService.HandleShortLink)))

	log.Fatal(rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		linkspb.RegisterLinksServiceServer(server, linksService)
	})))
}
