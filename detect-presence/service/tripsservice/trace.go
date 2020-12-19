package tripsservice

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/service/tripsservice"

var tracer = otel.Tracer(instrumentationName)
