package blocky

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	nomad "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "spx01/blocky"
	imageVersion = "sha256:59b3661951c28db0eecd9bb2e673c798d7c861d286e7713665da862e5254c477"
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
	job := nomadic.NewJob(a.name, 90)
	tg := nomadic.AddTaskGroup(job, "blocky", 1)
	nomadic.AddPort(tg, nomad.Port{Label: "dns"})
	nomadic.AddPort(tg, nomad.Port{Label: "http"})
	nomadic.AddService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "dns",
		TaskName:  "blocky",
		Tags:      []string{"dns"},
		Checks: []nomad.ServiceCheck{
			{
				Type:    "script",
				Command: "dig",
				Args: []string{
					"@${NOMAD_IP_dns}",
					"-p",
					"${NOMAD_HOST_PORT_dns}",
					"google.com",
				},
				Interval: 30 * time.Second,
				Timeout:  5 * time.Second,
			},
		},
	})
	nomadic.AddService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "http",
		Tags:      []string{"http"},
		Meta: map[string]string{
			"metrics_path": "/metrics",
		},
		Checks: []nomad.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/",
				Interval:             30 * time.Second,
				Timeout:              5 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "blocky",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"/app/blocky",
				"--config",
				"${NOMAD_TASK_DIR}/config.yaml",
			},
			"ports": []string{"dns", "http"},
		},
		Templates: []*nomad.Template{
			{
				EmbeddedTmpl: nomadic.String(configFile),
				DestPath:     nomadic.String("local/config.yaml"),
				ChangeMode:   nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(75),
		nomadic.WithLoggingTag(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}
	return nil
}

//go:embed config.yaml
var configFile string
