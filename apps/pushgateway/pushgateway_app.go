package pushgateway

import (
	"context"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "prom/pushgateway"
	imageVersion = "sha256:84327d5679194898b4952009b8f407e79a82f5f39dfbdfe8959bc5b62a84af0d"
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
	tg := nomadic.AddTaskGroup(job, "pushgateway", 1)

	nomadic.AddPort(tg, nomadapi.Port{Label: "http"})

	nomadic.AddService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "http",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/-/ready",
				Interval:             15 * time.Second,
				Timeout:              3 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "pushgateway",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"--web.listen-address=:${NOMAD_PORT_http}",
			},
			"ports": []string{"http"},
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
