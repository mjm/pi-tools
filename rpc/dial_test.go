package rpc_test

import (
	"context"
	"log"

	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/rpc"
)

func ExampleMustDial() {
	ctx := context.Background()

	tripsConn := rpc.MustDial(ctx, "localhost:2121")
	defer tripsConn.Close()

	trips := tripspb.NewTripsServiceClient(tripsConn)
	res, err := trips.ListTrips(ctx, &tripspb.ListTripsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("found %d trips", len(res.GetTrips()))
}
