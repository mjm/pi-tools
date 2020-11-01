package messagesservice

import (
	"go.opentelemetry.io/otel/api/global"
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/bot/service/messagesservice"

var tracer = global.Tracer(instrumentationName)
