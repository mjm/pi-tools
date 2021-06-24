package grafana

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "grafana/grafana"
	imageVersion = "sha256:f0817ecbf8dcf33e10cca2245bd25439433c441189bbe1ce935ac61d05f9cc6f"
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

	svcDefaults := &consulapi.ServiceConfigEntry{
		Kind:     consulapi.ServiceDefaults,
		Name:     a.name,
		Protocol: "http",
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", a.name, err)
	}

	svcIntentions := &consulapi.ServiceIntentionsConfigEntry{
		Kind: consulapi.ServiceIntentions,
		Name: a.name,
		Sources: []*consulapi.SourceIntention{
			{
				Name:       "ingress-http",
				Precedence: 9,
				Type:       consulapi.IntentionSourceConsul,
				Permissions: []*consulapi.IntentionPermission{
					{
						Action: consulapi.IntentionActionAllow,
						HTTP: &consulapi.IntentionHTTPPermission{
							PathPrefix: "/",
						},
					},
				},
			},
			{
				Action:     consulapi.IntentionActionDeny,
				Name:       "*",
				Precedence: 8,
				Type:       consulapi.IntentionSourceConsul,
			},
		},
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", err)
	}

	job := nomadic.NewJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "grafana", 3)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "3000",
		Checks: []nomadapi.ServiceCheck{
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

	nomadic.AddTask(tg, &nomadapi.Task{
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
	if err := clients.Vault.Sys().DeletePolicy("grafana"); err != nil {
		return fmt.Errorf("deleting grafana vault policy: %w", err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s job: %w", a.name, err)
	}

	return nil
}

var _ nomadic.Deployable = (*App)(nil)
