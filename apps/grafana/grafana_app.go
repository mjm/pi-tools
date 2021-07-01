package grafana

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
	imageRepo    = "grafana/grafana"
	// grafana 8.0.3
	imageVersion = "sha256:db7bdb09f965cc087bfa4aad75532f824a5eeef788602a54caf29e9d09be71a3"
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
		nomadic.AppAwareIntention("ingress-http",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/"))),
		nomadic.DenyIntention("*"))

	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", a.name, err)
	}

	job := nomadic.NewJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "grafana", 3)

	nomadic.AddConnectService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "3000",
		Checks: []nomad.ServiceCheck{
			{
				Type:     "http",
				Path:     "/api/health",
				Interval: 15 * time.Second,
				Timeout:  3 * time.Second,
			},
		},
	},
		nomadic.WithMetricsScraping("/metrics"),
		nomadic.WithUpstreams(
			nomadic.ConsulUpstream("loki", 3100)))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "grafana",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
		},
		Env: map[string]string{
			"GF_PATHS_CONFIG":       "${NOMAD_SECRETS_DIR}/grafana.ini",
			"GF_PATHS_PROVISIONING": "${NOMAD_TASK_DIR}/provisioning",
		},
		Templates: taskTemplates(),
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(100),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name))

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
