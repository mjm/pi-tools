package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/etherlabsio/healthcheck"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"

	"github.com/mjm/pi-tools/go-links/database/migrate"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
	"github.com/mjm/pi-tools/go-links/service/linksservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
	"github.com/mjm/pi-tools/storage"
)

func main() {
	rpc.SetDefaultHTTPPort(4240)
	rpc.SetDefaultGRPCPort(4241)
	storage.SetDefaultDBName("golinks_dev")
	flag.Parse()

	stopObs := observability.MustStart("go-links")
	defer stopObs()

	db := storage.MustOpenDB(migrate.Data)

	linksService := linksservice.New(db)
	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second),
			healthcheck.WithChecker("database", db))))
	http.Handle("/",
		otelhttp.WithRouteTag("HandleShortLink", http.HandlerFunc(linksService.HandleShortLink)))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		linkspb.RegisterLinksServiceServer(server, linksService)
	}))

	signal.Wait()
}
