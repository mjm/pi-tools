package unifiexporter

import (
	"context"
	_ "embed"
	"fmt"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const imageRepo = "mmoriarity/unifi_exporter"

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
		return fmt.Errorf("updating %s vault policy: %w", a.name, err)
	}

	job := nomadic.NewJob(a.name, 30)
	tg := nomadic.AddTaskGroup(job, "unifi-exporter", 1)
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "http",
		To:    9130,
	})

	nomadic.AddService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "http",
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "unifi-exporter",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
			"command": "/unifi_exporter",
			"args": []string{
				"-config.file=${NOMAD_SECRETS_DIR}/config.yml",
			},
			"ports": []string{"http"},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: &configFile,
				DestPath:     nomadic.String("secrets/config.yml"),
				ChangeMode:   nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(30),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	panic("implement me")
}

//go:embed unifi-exporter.hcl
var vaultPolicy string

//go:embed config.yml
var configFile string
