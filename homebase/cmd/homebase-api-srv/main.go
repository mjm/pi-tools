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

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/service/apiservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	tripsURL   = flag.String("trips-url", "localhost:2121", "URL for trips service")
	deployURL  = flag.String("deploy-url", "localhost:8481", "URL for deploy service")
	schemaPath = flag.String("schema-path", "/schema.graphql", "Path to the file with the GraphQL schema")
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

	deployConn := rpc.MustDial(ctx, *deployURL)
	defer deployConn.Close()

	deploy := deploypb.NewDeployServiceClient(deployConn)

	apiService, err := apiservice.New(string(schema), trips, deploy)
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
