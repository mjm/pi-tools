package nomadic

import (
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/mjm/pi-tools/pkg/nomadic")
