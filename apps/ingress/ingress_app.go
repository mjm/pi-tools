package ingress

import (
	"context"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "nginx"
	imageVersion = "sha256:763d95e3db66d9bd1bb926c029e5659ee67eb49ff57f83d331de5f5af6d2ae0c"
)

type App struct {
	name string
}

var _ nomadic.Deployable = (*App)(nil)

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

	httpServiceName := a.name + "-http"
	httpsServiceName := a.name + "-https"

	svcDefaults := &consulapi.ServiceConfigEntry{
		Kind:     consulapi.ServiceDefaults,
		Name:     httpServiceName,
		Protocol: "http",
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("updating %s service defaults: %w", httpServiceName, err)
	}

	svcIntentions := &consulapi.ServiceIntentionsConfigEntry{
		Kind: consulapi.ServiceIntentions,
		Name: httpServiceName,
		Sources: []*consulapi.SourceIntention{
			{
				Action:     consulapi.IntentionActionDeny,
				Name:       "*",
				Precedence: 8,
				Type:       consulapi.IntentionSourceConsul,
			},
		},
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("updating %s service intentions: %w", httpServiceName, err)
	}

	job := nomadic.NewJob(a.name, 70)
	job.Update = &nomadapi.UpdateStrategy{
		MaxParallel: nomadic.Int(1),
		Stagger:     nomadic.Duration(30 * time.Second),
	}

	tg := nomadic.AddTaskGroup(job, "ingress", 2)
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "http",
		To:    80,
		Value: 80,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "https",
		To:    443,
		Value: 443,
	})

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      httpServiceName,
		PortLabel: "80",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:     "http",
				Path:     "/healthz",
				Interval: 15 * time.Second,
				Timeout:  3 * time.Second,
				Header: map[string][]string{
					// Just need to have some Host header
					"Host": {"ingress"},
				},
			},
		},
	}, nomadic.WithUpstreams(connectUpstreams()...))
	nomadic.AddService(tg, &nomadapi.Service{
		Name:      httpsServiceName,
		PortLabel: "443",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "nginx",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"volumes": []string{
				"local:/etc/nginx/conf.d",
				"secrets:/etc/nginx/ssl",
			},
		},
		Templates: taskTemplates(),
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name),
		nomadic.WithVaultChangeNoop())

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return err
	}
	return nil
}
