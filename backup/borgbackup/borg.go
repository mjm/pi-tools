package borgbackup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

type Borg struct {
	borgPath string
}

func New(borgPath string) *Borg {
	return &Borg{
		borgPath: borgPath,
	}
}

func (b *Borg) commandJSON(ctx context.Context, result interface{}, args ...string) error {
	ctx, span := tracer.Start(ctx, "BorgCommand",
		trace.WithAttributes(
			label.String("borg.path", b.borgPath),
			label.String("borg.args", strings.Join(args, " "))))
	defer span.End()

	var realArgs []string
	realArgs = append(realArgs, args...)
	realArgs = append(realArgs, "--json")
	out, err := exec.CommandContext(ctx, b.borgPath, realArgs...).Output()
	if err != nil {
		span.RecordError(err)

		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return spanerr.RecordError(ctx, fmt.Errorf("borg: %s", exitError.Stderr))
		}
		return spanerr.RecordError(ctx, err)
	}

	return json.Unmarshal(out, result)
}
