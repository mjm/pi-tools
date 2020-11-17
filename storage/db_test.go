package storage_test

import (
	"context"
	"flag"
	"log"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/database/migrate"
	"github.com/mjm/pi-tools/storage"
)

func ExampleMustOpenDB() {
	storage.SetDefaultDBName("trips_dev")
	flag.Parse()

	db := storage.MustOpenDB(migrate.Data)
	q := database.New(db)

	ctx := context.Background()
	trips, err := q.ListTrips(ctx, 30)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("got %d trips", len(trips))
}
