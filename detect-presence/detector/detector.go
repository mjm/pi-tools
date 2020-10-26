package detector

import (
	"context"

	"go.opentelemetry.io/otel/api/global"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/detector"

var tracer = global.Tracer(instrumentationName)

type Detector interface {
	IsHealthy(ctx context.Context) (bool, error)
	DetectDevice(ctx context.Context, addr string) (bool, error)
}
