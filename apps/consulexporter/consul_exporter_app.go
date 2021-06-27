package consulexporter

import (
	"context"
	_ "embed"
	"fmt"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "prom/consul-exporter"
	imageVersion = "sha256:4e45d018f2fd35afbc3c0c79aa6fe9f43642f9fe49170aca989998015c76c922"
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
		return fmt.Errorf("updating %s vault policy: %w", a.name, err)
	}

	job := nomadic.NewJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "consul-exporter", 1)

	nomadic.AddPort(tg, nomadapi.Port{Label: "http"})
	nomadic.AddService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "http",
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "consul-exporter",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"--web.listen-address=:${NOMAD_PORT_http}",
				"--consul.server=${attr.unique.network.ip-address}:8500",
			},
			"ports": []string{
				"http",
			},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: nomadic.String(`
{{ with secret "consul/creds/prometheus" }}
CONSUL_HTTP_TOKEN={{ .Data.token }}
{{ end }}
`),
				DestPath:   nomadic.String("secrets/consul.env"),
				Envvars:    nomadic.Bool(true),
				ChangeMode: nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name, err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}

//go:embed consul-exporter.hcl
var vaultPolicy string
