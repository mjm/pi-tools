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
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
	"github.com/mjm/pi-tools/go-links/service/linksservice"
)

var (
	httpPort = flag.Int("http-port", 4240, "HTTP port to listen on for metrics and API requests")
	dbDSN    = flag.String("db", "dbname=golinks_dev sslmode=disable", "Connection string for connecting to PostgreSQL database for storing links")
	debug    = flag.Bool("debug", false, "Show debug tracing output in stdout")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if *debug {
		exp, err := stdout.NewExporter(stdout.WithPrettyPrint())
		if err != nil {
			log.Panicf("failed to initialize stdout exporter %v\n", err)
		}
		bsp := sdktrace.NewBatchSpanProcessor(exp)
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithConfig(
				sdktrace.Config{
					DefaultSampler: sdktrace.AlwaysSample(),
				},
			),
			sdktrace.WithSpanProcessor(bsp),
		)
		global.SetTracerProvider(tp)
	}

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

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), handler))
}
