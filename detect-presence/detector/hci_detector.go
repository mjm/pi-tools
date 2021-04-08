package detector

import (
	"context"
	"io/ioutil"
	"os/exec"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type HCIDetector struct {
	DeviceName string
}

func (d *HCIDetector) IsHealthy(ctx context.Context) (bool, error) {
	ctx, span := tracer.Start(ctx, "HCIDetector.IsHealthy",
		trace.WithAttributes(attribute.String("bluetooth.device.name", d.DeviceName)))
	defer span.End()

	cmd := exec.CommandContext(ctx, "/bin/hciconfig", d.DeviceName)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return false, err
	}

	if err := cmd.Start(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	output, err := ioutil.ReadAll(outPipe)
	if err != nil {
		return false, err
	}
	errOutput, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return false, err
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			span.SetAttributes(attribute.Int("cmd.exit_code", exitErr.ExitCode()))

			if exitErr.ExitCode() == 1 {
				span.SetStatus(codes.Error, string(errOutput))
				return false, nil
			}
		}
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	healthy := strings.Contains(string(output), "UP RUNNING")
	span.SetAttributes(attribute.Bool("bluetooth.healthy", healthy))
	return healthy, nil
}

func (*HCIDetector) DetectDevice(ctx context.Context, addr string) (bool, error) {
	ctx, span := tracer.Start(ctx, "HCIDetector.DetectDevice",
		trace.WithAttributes(attribute.String("device.addr", addr)))
	defer span.End()

	cmd := exec.CommandContext(ctx, "/usr/bin/hcitool", "info", addr)
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return false, err
	}

	if err := cmd.Start(); err != nil {
		return false, err
	}

	errOutput, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return false, err
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			span.SetAttributes(attribute.Int("cmd.exit_code", exitErr.ExitCode()))

			if exitErr.ExitCode() == 1 {
				span.SetAttributes(attribute.Bool("device.present", false))
				span.SetStatus(codes.Error, string(errOutput))
				return false, nil
			}
		}
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	span.SetAttributes(attribute.Bool("device.present", true))
	return true, nil
}
