package nut

import (
	"context"
	"fmt"
	"strings"
	"time"

	nomad "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo = "instantlinux/nut-upsd"
	// nut 2.7.4-r8
	imageVersion = "sha256:f4948378c7e7c27c857e822aac174854389b8607fbe937538ef647d9cade69b0"

	exporterImageRepo = "mmoriarity/nut_exporter"
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
	job := nomadic.NewJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "nut", 1)

	nomadic.AddPort(tg, nomad.Port{
		Label: "metrics",
		To:    9199,
	})

	nomadic.AddService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "metrics",
		Checks: []nomad.ServiceCheck{
			{
				Type:     "http",
				Path:     "/",
				Interval: 30 * time.Second,
				Timeout:  5 * time.Second,
			},
		},
	}, nomadic.WithMetricsScraping("/ups_metrics"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "nut-exporter",
		Config: map[string]interface{}{
			"image":   nomadic.Image(exporterImageRepo, "latest"),
			"command": "/nut_exporter",
			"args": []string{
				"--nut.server=10.0.0.10",
				"--nut.vars_enable=" + strings.Join(enabledVariables, ","),
			},
			"ports": []string{"metrics"},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}
