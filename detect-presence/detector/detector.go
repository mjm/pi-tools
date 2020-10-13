package detector

import (
	"context"
)

type Detector interface {
	IsHealthy(ctx context.Context, deviceName string) (bool, error)
	DetectDevice(ctx context.Context, addr string) (bool, error)
}
