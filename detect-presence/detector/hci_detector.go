package detector

import (
	"context"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type HCIDetector struct{}

func (*HCIDetector) IsHealthy(ctx context.Context, deviceName string) (bool, error) {
	cmd := exec.CommandContext(ctx, "/bin/hciconfig", deviceName)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return false, err
	}

	if err := cmd.Start(); err != nil {
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
			if exitErr.ExitCode() == 1 {
				log.Printf("Detection failed, output: %s", errOutput)
				return false, nil
			}
		}
		return false, err
	}

	if strings.Contains(string(output), "UP RUNNING") {
		return true, nil
	}

	return false, nil
}

func (*HCIDetector) DetectDevice(ctx context.Context, addr string) (bool, error) {
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
			if exitErr.ExitCode() == 1 {
				log.Printf("Detection failed, output: %s", errOutput)
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}
