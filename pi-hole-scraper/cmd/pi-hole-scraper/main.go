package main

import (
	"context"
	"flag"
	"log"

	"github.com/mjm/pi-tools/pi-hole-scraper/scraper"
)

var (
	piHoleAddress  = flag.String("pi-hole-address", "http://localhost", "Address of the PiHole server to scrape from")
	piHolePassword = flag.String("pi-hole-password", "", "Password for accessing PiHole API (WEBPASSWORD in setupVars.conf)")
)

func main() {
	flag.Parse()

	s := scraper.New(*piHoleAddress, *piHolePassword)
	ctx := context.Background()

	stats, err := s.GetStats(ctx)
	if err != nil {
		log.Fatalf("Error scraping stats from PiHole: %v", err)
	}

	log.Printf("Stats: %#v", stats)
}
