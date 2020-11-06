package rpc_test

import (
	"flag"

	"google.golang.org/grpc"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

func ExampleWithRegisteredServices() {
	rpc.SetDefaultHTTPPort(1234)
	flag.Parse()

	var tripsSrv tripspb.TripsServiceServer
	rpc.ListenAndServe(rpc.WithRegisteredServices(func(server *grpc.Server) {
		tripspb.RegisterTripsServiceServer(server, tripsSrv)
	}))

	signal.Wait()
}
