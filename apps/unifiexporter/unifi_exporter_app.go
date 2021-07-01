package unifiexporter

import (
	"context"
	_ "embed"
	"fmt"

	nomad "github.com/hashicorp/nomad/api"

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
	nomadic.AddPort(tg, nomad.Port{
		Label: "http",
		To:    9130,
	})

	nomadic.AddService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "http",
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "unifi-exporter",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
			"command": "/unifi_exporter",
			"args": []string{
				"-config.file=${NOMAD_SECRETS_DIR}/config.yml",
			},
			"ports": []string{"http"},
		},
		Templates: []*nomad.Template{
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
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name, err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s job: %w", a.name, err)
	}

	return nil
}

//go:embed unifi-exporter.hcl
var vaultPolicy string

//go:embed config.yml
var configFile string
