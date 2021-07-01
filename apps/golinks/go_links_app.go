package golinks

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
	imageRepo = "mmoriarity/go-links-srv"
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
		nomadic.AllowIntention("ingress-http"),
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", a.name, err)
	}

	grpcName := a.name + "-grpc"
	svcDefaults = nomadic.NewServiceDefaults(grpcName, "grpc")
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", grpcName, err)
	}

	svcIntentions = nomadic.NewServiceIntentions(grpcName,
		nomadic.AppAwareIntention("homebase-api",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/LinksService/")),
			nomadic.DenyAllHTTP()),
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", grpcName, err)
	}

	job := nomadic.NewJob(a.name, 50)
	tg := nomadic.AddTaskGroup(job, "go-links", 2)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "4240",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/healthz",
				Interval:             15 * time.Second,
				Timeout:              3 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      grpcName,
		PortLabel: "4241",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "go-links-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
			"command": "/go-links",
			"args": []string{
				"-db",
				"dbname=golinks host=10.0.2.102 sslmode=disable",
			},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: nomadic.String(`
{{ with secret "database/creds/go-links" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
`),
				DestPath:   nomadic.String("secrets/db.env"),
				Envvars:    nomadic.Bool(true),
				ChangeMode: nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(50),
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

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", a.name, err)
	}
	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", a.name, err)
	}

	grpcName := a.name + "-grpc"
	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, grpcName, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", grpcName, err)
	}
	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, grpcName, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", grpcName, err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}

//go:embed go-links.hcl
var vaultPolicy string
