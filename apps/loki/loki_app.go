package loki

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo = "grafana/loki"
	// loki 2.2.1
	imageVersion = "sha256:7d2ddbe46c11cf9778eba0abf67bc963366dcfd7bda1a123e5244187e64dafec"
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
	if err := clients.Vault.Sys().PutPolicy(a.name, vaultPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", err)
	}

	job := nomadic.NewJob(a.name, 80)
	tg := nomadic.AddTaskGroup(job, "loki", 1)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "3100",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:     "http",
				Path:     "/ready",
				Interval: 15 * time.Second,
				Timeout:  3 * time.Second,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "loki",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"-config.file=${NOMAD_TASK_DIR}/loki.yml",
			},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: nomadic.String(configFile),
				DestPath:     nomadic.String("local/loki.yml"),
				ChangeMode:   nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(150),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name),
		nomadic.WithVaultChangeNoop())

	resp, _, err := clients.Nomad.Jobs().Plan(job, true, nil)
	if err != nil {
		return fmt.Errorf("planning %s job: %w", *job.ID, err)
	}
	if resp.Diff.Type == "None" {
		return nil
	}

	if _, _, err := clients.Nomad.Jobs().Register(job, nil); err != nil {
		return fmt.Errorf("registering %s job: %w", *job.ID, err)
	}
	return nil
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s job: %w", a.name, err)
	}

	return nil
}

var _ nomadic.Deployable = (*App)(nil)

//go:embed loki.hcl
var vaultPolicy string

//go:embed loki.yml
var configFile string
