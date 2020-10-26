package debug

import (
	"flag"
)

var debug = flag.Bool("debug", false, "Run with support for easier debugging")

func IsEnabled() bool {
	return *debug
}
