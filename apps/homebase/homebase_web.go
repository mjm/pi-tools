package homebase

import (
	"context"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	webImageRepo = "mmoriarity/homebase-srv"
)

func (a *App) addWebTaskGroup(job *nomadapi.Job) {
	tg := nomadic.AddTaskGroup(job, "homebase-srv", 2)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "3000",
	}, nomadic.WithUpstreams(
		nomadic.ConsulUpstream("homebase-api", 6460)))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "homebase-srv",
		Config: map[string]interface{}{
			"image": nomadic.Image(webImageRepo, "latest"),
		},
		Env: map[string]string{
			"GRAPHQL_URL": "http://localhost:6460/graphql",
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(75),
		nomadic.WithLoggingTag(a.name))
}

func (a *App) installWebConfigEntries(ctx context.Context, clients nomadic.Clients) error {
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

	return nil
}

func (a *App) uninstallWebConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", a.name, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", a.name, err)
	}

	return nil
}
