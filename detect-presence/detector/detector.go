package detector

import (
	"context"
)

type Detector interface {
	IsHealthy(ctx context.Context) (bool, error)
	DetectDevice(ctx context.Context, addr string) (bool, error)
}
