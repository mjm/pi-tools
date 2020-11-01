package telegram

import (
	"go.opentelemetry.io/otel/api/global"
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/bot/telegram"

var tracer = global.Tracer(instrumentationName)
