package appservice

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/service/appservice"

var tracer = otel.Tracer(instrumentationName)
