package homebase

import (
	"context"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	apiImageRepo = "mmoriarity/homebase-api-srv"
)

func (a *App) addAPITaskGroup(job *nomadapi.Job) {
	tg := nomadic.AddTaskGroup(job, "homebase-api", 2)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name + "-api",
		PortLabel: "6460",
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
			nomadic.ConsulUpstream("go-links-grpc", 4241),
			nomadic.ConsulUpstream("detect-presence-grpc", 2121),
			nomadic.ConsulUpstream("deploy-grpc", 8481),
			nomadic.ConsulUpstream("backup-grpc", 2321)))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "homebase-api-srv",
		Config: map[string]interface{}{
			"image": nomadic.Image(apiImageRepo, "latest"),
			"command": "/homebase-api-srv",
			"args": []string{
				"-prometheus-url",
				"http://prometheus.service.consul:9090",
				"-paperless-url",
				"http://paperless.service.consul:8000",
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name + "-api"),
		nomadic.WithTracingEnv())
}

func (a *App) installAPIConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	name := a.name + "-api"

	svcDefaults := &consulapi.ServiceConfigEntry{
		Kind:     consulapi.ServiceDefaults,
		Name:     name,
		Protocol: "http",
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", name, err)
	}

	svcIntentions := &consulapi.ServiceIntentionsConfigEntry{
		Kind: consulapi.ServiceIntentions,
		Name: name,
		Sources: []*consulapi.SourceIntention{
			{
				Name:       "ingress-http",
				Precedence: 9,
				Type:       consulapi.IntentionSourceConsul,
				Permissions: []*consulapi.IntentionPermission{
					{
						Action: consulapi.IntentionActionAllow,
						HTTP: &consulapi.IntentionHTTPPermission{
							PathExact: "/graphql",
						},
					},
					{
						Action: consulapi.IntentionActionDeny,
						HTTP: &consulapi.IntentionHTTPPermission{
							PathPrefix: "/",
						},
					},
				},
			},
			{
				Name:       a.name,
				Precedence: 9,
				Type:       consulapi.IntentionSourceConsul,
				Permissions: []*consulapi.IntentionPermission{
					{
						Action: consulapi.IntentionActionAllow,
						HTTP: &consulapi.IntentionHTTPPermission{
							PathExact: "/graphql",
						},
					},
					{
						Action: consulapi.IntentionActionDeny,
						HTTP: &consulapi.IntentionHTTPPermission{
							PathPrefix: "/",
						},
					},
				},
			},
			{
				Name:       "*",
				Precedence: 8,
				Type:       consulapi.IntentionSourceConsul,
				Action:     consulapi.IntentionActionDeny,
			},
		},
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", name, err)
	}

	return nil
}

func (a *App) uninstallAPIConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	name := a.name + "-api"

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, name, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", name, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, name, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", name, err)
	}

	return nil
}
