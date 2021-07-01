package trippliteexporter

import (
	"context"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo = "mmoriarity/tripplite-exporter"
)

type App struct {
	name string
}

func New(name string) *App {
	return &App{
		name: name,
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Install(ctx context.Context, clients nomadic.Clients) error {
	job := nomadic.NewSystemJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "tripplite-exporter", 1)

	// this needs to connect to the Tripplite UPS over USB,
	// and it's only plugged in to this machine.
	job.Constrain(&nomadapi.Constraint{
		LTarget: "${node.unique.name}",
		Operand: "=",
		RTarget: "raspberrypi",
	})

	nomadic.AddPort(tg, nomadapi.Port{
		Label: "http",
		To:    8080,
	})

	nomadic.AddService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "http",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/healthz",
				Interval:             30 * time.Second,
				Timeout:              5 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "tripplite-exporter",
		Config: map[string]interface{}{
			"image":      nomadic.Image(imageRepo, "latest"),
			"command":    "/tripplite_exporter",
			"ports":      []string{"http"},
			"privileged": true,
			"mount": []map[string]interface{}{
				{
					"type":   "bind",
					"target": "/dev/bus/usb",
					"source": "/dev/bus/usb",
				},
			},
		},
	},
		nomadic.WithCPU(30),
		nomadic.WithMemoryMB(30),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithTracingEnv())

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}
