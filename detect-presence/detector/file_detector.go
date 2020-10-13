package detector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type FileDetector struct {
	Path string
}

func (d *FileDetector) IsHealthy(ctx context.Context) (bool, error) {
	return true, nil
}

func (d *FileDetector) DetectDevice(ctx context.Context, addr string) (bool, error) {
	f, err := os.Open(d.Path)
	if err != nil {
		return false, fmt.Errorf("opening detector file: %w", err)
	}
	defer f.Close()

	var devices map[string]bool
	if err := json.NewDecoder(f).Decode(&devices); err != nil {
		return false, err
	}

	present, ok := devices[addr]
	if !ok {
		return false, fmt.Errorf("no such device %q", addr)
	}

	return present, nil
}
