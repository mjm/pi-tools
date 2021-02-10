package borgbackup

import (
	"go.opentelemetry.io/otel"
)

const instrumentationName = "github.com/mjm/pi-tools/backup/borgbackup"

var tracer = otel.Tracer(instrumentationName)
