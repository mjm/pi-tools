package messagesservice

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/bot/service/messagesservice"

var tracer = otel.Tracer(instrumentationName)
