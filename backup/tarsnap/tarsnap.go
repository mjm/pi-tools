package tarsnap

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type Tarsnap struct {
	tarsnapPath string
}

func New(path string) *Tarsnap {
	return &Tarsnap{
		tarsnapPath: path,
	}
}

func (t *Tarsnap) runCommand(ctx context.Context, args ...string) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "TarsnapCommand",
		trace.WithAttributes(
			label.String("tarsnap.path", t.tarsnapPath),
			label.String("tarsnap.args", strings.Join(args, " "))))
	defer span.End()

	out, err := exec.CommandContext(ctx, t.tarsnapPath, args...).Output()
	if err != nil {
		span.RecordError(err)

		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return nil, spanerr.RecordError(ctx, fmt.Errorf("tarsnap: %s", exitError.Stderr))
		}
		return nil, spanerr.RecordError(ctx, err)
	}

	return out, nil
}
