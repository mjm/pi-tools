package tarsnap

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/backup/tarsnap"

var tracer = otel.Tracer(instrumentationName)
