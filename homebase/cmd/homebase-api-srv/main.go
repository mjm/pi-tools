package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
	"github.com/mjm/pi-tools/homebase/service/apiservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	tripsURL      = flag.String("trips-url", "127.0.0.1:2121", "URL for trips service")
	linksURL      = flag.String("links-url", "127.0.0.1:4241", "URL for links service")
	deployURL     = flag.String("deploy-url", "127.0.0.1:8481", "URL for deploy service")
	backupURL     = flag.String("backup-url", "127.0.0.1:2321", "URL for backup service")
	prometheusURL = flag.String("prometheus-url", "https://prometheus.home.mattmoriarity.com", "URL for Prometheus for querying alerts")
	paperlessURL  = flag.String("paperless-url", "https://paperless.home.mattmoriarity.com", "URL for paperless-ng")
	schemaPath    = flag.String("schema-path", "/schema.graphql", "Path to the file with the GraphQL schema")
)

func main() {
	rpc.SetDefaultHTTPPort(6460)
	flag.Parse()

	stopObs := observability.MustStart("homebase-api-srv")
	defer stopObs()

	ctx := context.Background()

	schema, err := ioutil.ReadFile(*schemaPath)
	if err != nil {
		log.Panicf("reading GraphQL schema: %v", err)
	}

	tripsConn := rpc.MustDial(ctx, *tripsURL)
	defer tripsConn.Close()

	trips := tripspb.NewTripsServiceClient(tripsConn)

	linksConn := rpc.MustDial(ctx, *linksURL)
	defer linksConn.Close()

	links := linkspb.NewLinksServiceClient(linksConn)

	deployConn := rpc.MustDial(ctx, *deployURL)
	defer deployConn.Close()

	deploy := deploypb.NewDeployServiceClient(deployConn)

	backupConn := rpc.MustDial(ctx, *backupURL)
	defer backupConn.Close()

	backup := backuppb.NewBackupServiceClient(backupConn)

	apiService, err := apiservice.New(string(schema), trips, links, deploy, backup, *prometheusURL, *paperlessURL)
	if err != nil {
		log.Panicf("creating API service: %v", err)
	}

	http.Handle("/graphql", apiService)
	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))

	go rpc.ListenAndServe()

	signal.Wait()
}
