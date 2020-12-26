package deployservice

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/deploy/service/deployservice"

var tracer = otel.Tracer(instrumentationName)
