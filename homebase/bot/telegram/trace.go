package telegram

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/homebase/bot/telegram"

var tracer = otel.Tracer(instrumentationName)
