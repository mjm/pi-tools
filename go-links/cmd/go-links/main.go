package main

import (
	"flag"
)

var (
	httpPort = flag.Int("http-port", 4240, "HTTP port to listen on for metrics and API requests")
	dbDSN    = flag.String("db", ":memory:", "Connection string for connecting to SQLite3 database for storing links")
)

func main() {
	flag.Parse()
}
