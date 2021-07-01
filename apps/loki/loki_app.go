package loki

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"

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
		return fmt.Errorf("updating %s vault policy: %w", a.name, err)
	}

	svcDefaults := nomadic.NewServiceDefaults(a.name, "http")
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", a.name, err)
	}

	svcIntentions := nomadic.NewServiceIntentions(a.name,
		nomadic.AppAwareIntention("grafana",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/"))),
		nomadic.AppAwareIntention("promtail",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/"))),
		nomadic.DenyIntention("*"))

	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", a.name, err)
	}

	job := nomadic.NewJob(a.name, 80)
	tg := nomadic.AddTaskGroup(job, "loki", 1)

	nomadic.AddConnectService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "3100",
		Checks: []nomad.ServiceCheck{
			{
				Type:     "http",
				Path:     "/ready",
				Interval: 15 * time.Second,
				Timeout:  3 * time.Second,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "loki",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"-config.file=${NOMAD_TASK_DIR}/loki.yml",
			},
		},
		Templates: []*nomad.Template{
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

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consul.ServiceDefaults, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", a.name, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consul.ServiceIntentions, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", a.name, err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s job: %w", a.name, err)
	}

	return nil
}

//go:embed loki.hcl
var vaultPolicy string

//go:embed loki.yml
var configFile string
