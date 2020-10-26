package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
	"github.com/mjm/pi-tools/go-links/service/linksservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/instrumentation/otelsql"
)

var (
	httpPort = flag.Int("http-port", 4240, "HTTP port to listen on for metrics and API requests")
	dbDSN    = flag.String("db", "dbname=golinks_dev sslmode=disable", "Connection string for connecting to PostgreSQL database for storing links")
)

func main() {
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
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	linkspb.RegisterLinksServiceServer(grpcServer, linksService)
	reflection.Register(grpcServer)

	wrappedGrpc := grpcweb.WrapServer(grpcServer, grpcweb.WithOriginFunc(func(origin string) bool {
		if origin == "http://localhost:8080" || origin == "http://mars.local:8080" {
			return true
		}
		log.Printf("Rejecting unknown origin: %s", origin)
		return false
	}))

	handler := otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wrappedGrpc.IsAcceptableGrpcCorsRequest(r) || wrappedGrpc.IsGrpcWebRequest(r) {
			wrappedGrpc.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			log.Print(r)
			grpcServer.ServeHTTP(w, r)
			return
		}

		http.DefaultServeMux.ServeHTTP(w, r)
	}), "Server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))

	http.Handle("/", otelhttp.WithRouteTag("HandleShortLink", http.HandlerFunc(linksService.HandleShortLink)))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), handler))
}
