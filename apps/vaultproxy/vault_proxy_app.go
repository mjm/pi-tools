package vaultproxy

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
	imageRepo = "mmoriarity/vault-proxy"
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

	svcDefaults := &consul.ServiceConfigEntry{
		Kind:     consul.ServiceDefaults,
		Name:     a.name,
		Protocol: "http",
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", a.name, err)
	}

	svcIntentions := &consul.ServiceIntentionsConfigEntry{
		Kind: consul.ServiceIntentions,
		Name: a.name,
		Sources: []*consul.SourceIntention{
			{
				Name:       "ingress-http",
				Precedence: 9,
				Type:       consul.IntentionSourceConsul,
				Permissions: []*consul.IntentionPermission{
					{
						Action: consul.IntentionActionAllow,
						HTTP: &consul.IntentionHTTPPermission{
							PathPrefix: "/",
						},
					},
				},
			},
			{
				Name:       "*",
				Precedence: 8,
				Type:       consul.IntentionSourceConsul,
				Action:     consul.IntentionActionDeny,
			},
		},
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", a.name, err)
	}

	job := nomadic.NewJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "vault-proxy", 2)

	nomadic.AddConnectService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "2220",
		Checks: []nomad.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/healthz",
				Interval:             15 * time.Second,
				Timeout:              3 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "vault-proxy",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
			"command": "/vault-proxy",
		},
		Env: map[string]string{
			"VAULT_ADDR": "http://active.vault.service.consul:8200",
		},
		Templates: []*nomad.Template{
			{
				EmbeddedTmpl: nomadic.String(`
{{ with secret "kv/vault-proxy" }}
COOKIE_KEY={{ .Data.data.cookie_secret }}
{{ end }}
`),
				DestPath:   nomadic.String("secrets/proxy.env"),
				Envvars:    nomadic.Bool(true),
				ChangeMode: nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name),
		nomadic.WithTracingEnv())

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
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}

//go:embed vault-proxy.hcl
var vaultPolicy string
