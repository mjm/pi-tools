package linksservice

import (
	"context"
	"log"
	"os"
	"testing"

	"zombiezen.com/go/postgrestest"
)

var dbSrv *postgrestest.Server

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	dbSrv, err = postgrestest.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	exitCode := m.Run()
	dbSrv.Cleanup()
	os.Exit(exitCode)
}
