package main

import (
	"flag"
	"log"
	"net/http"

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
	storage.SetDefaultDBName("golinks_dev")
	flag.Parse()

	stopObs, err := observability.Start("go-links")
	if err != nil {
		log.Panicf("Error setting up observability: %v", err)
	}
	defer stopObs()

	db, err := storage.OpenDB(migrate.Data)
	if err != nil {
		log.Panicf("Error setting up storage: %v", err)
	}

	linksService := linksservice.New(db)
	http.Handle("/", otelhttp.WithRouteTag("HandleShortLink", http.HandlerFunc(linksService.HandleShortLink)))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		linkspb.RegisterLinksServiceServer(server, linksService)
	}))

	signal.Wait()
}
