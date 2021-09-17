package homebase

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
	botImageRepo = "mmoriarity/homebase-bot-srv"
)

func (a *App) addBotTaskGroup(job *nomadapi.Job) {
	tg := nomadic.AddTaskGroup(job, "homebase-bot", 2)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name + "-bot",
		PortLabel: "6360",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/healthz",
				Interval:             15 * time.Second,
				Timeout:              3 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	},
		nomadic.WithMetricsScraping("/metrics"),
		nomadic.WithUpstreams(
			nomadic.ConsulUpstream("detect-presence-grpc", 2121)))

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name + "-bot-grpc",
		PortLabel: "6361",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "homebase-bot-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(botImageRepo, "latest"),
			"command": "/homebase-bot-srv",
			"args": []string{
				"-db",
				"dbname=homebase_bot host=10.0.2.102 sslmode=disable",
				"-leader-elect",
			},
		},
		Env: map[string]string{
			"CONSUL_HTTP_ADDR": "${attr.unique.network.ip-address}:8500",
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: nomadic.String(`
{{ with secret "database/creds/homebase-bot" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
{{ with secret "kv/homebase-bot" }}
TELEGRAM_TOKEN={{ .Data.data.telegram_token | toJSON }}
{{ end }}
`),
				DestPath:   nomadic.String("secrets/secrets.env"),
				Envvars:    nomadic.Bool(true),
				ChangeMode: nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name+"-bot"),
		nomadic.WithVaultPolicies(a.name+"-bot"),
		nomadic.WithTracingEnv())
}

func (a *App) installBotVaultPolicy(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().PutPolicy(a.name+"-bot", botVaultPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", a.name+"-bot", err)
	}

	return nil
}

func (a *App) installBotConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	httpName := a.name + "-bot"

	svcDefaults := nomadic.NewServiceDefaults(httpName, "http")
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", httpName, err)
	}

	svcIntentions := nomadic.NewServiceIntentions(httpName,
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", httpName, err)
	}

	grpcName := httpName + "-grpc"

	svcDefaults = nomadic.NewServiceDefaults(grpcName, "grpc")
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", grpcName, err)
	}

	svcIntentions = nomadic.NewServiceIntentions(grpcName,
		nomadic.AppAwareIntention("detect-presence",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/MessagesService/")),
			nomadic.DenyAllHTTP()),
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", grpcName, err)
	}

	return nil
}

func (a *App) uninstallBotVaultPolicy(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name + "-bot"); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name+"-bot", err)
	}

	return nil
}

func (a *App) uninstallBotConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	httpName := a.name + "-bot"

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, httpName, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", httpName, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, httpName, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", httpName, err)
	}

	grpcName := httpName + "-grpc"

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, grpcName, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", grpcName, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, grpcName, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", grpcName, err)
	}

	return nil
}

//go:embed homebase-bot.hcl
var botVaultPolicy string
