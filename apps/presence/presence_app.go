package presence

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
	imageRepo       = "mmoriarity/detect-presence-srv"
	beaconImageRepo = "mmoriarity/beacon-srv"
)

type App struct {
	name       string
	beaconName string
}

func New(name, beaconName string) *App {
	return &App{
		name:       name,
		beaconName: beaconName,
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Install(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().PutPolicy(a.name, vaultPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", a.name, err)
	}

	if err := a.installConsulConfigEntries(ctx, clients); err != nil {
		return err
	}

	job := nomadic.NewJob(a.name, 50)
	tg := nomadic.AddTaskGroup(job, "detect-presence", 2)

	nomadic.AddConnectService(tg, &nomadapi.Service{
		// TODO work out the names so this can use a.name
		Name:      "detect-presence",
		PortLabel: "2120",
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
			nomadic.ConsulUpstream("homebase-bot-grpc", 6361)))

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      "detect-presence-grpc",
		PortLabel: "2121",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "detect-presence-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
			"command": "/detect-presence-srv",
			"args": []string{
				"-db",
				"dbname=presence host=10.0.2.102 sslmode=disable",
				"-mode",
				"client",
			},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: nomadic.String(`{{ with secret "kv/deploy" }}{{ .Data.data.github_token }}{{ end }}`),
				DestPath:     nomadic.String("secrets/github-token"),
				ChangeMode:   nomadic.String("restart"),
			},
			{
				EmbeddedTmpl: nomadic.String(`
{{ with secret "database/creds/presence" }}
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

	beaconJob := nomadic.NewSystemJob(a.beaconName, 50)
	tg = nomadic.AddTaskGroup(beaconJob, "beacon", 1)

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "beacon-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(beaconImageRepo, "latest"),
			"command": "/beacon-srv",
			"args": []string{
				"-proximity-uuid",
				"7298c12b-f658-445f-b1f2-5d6d582f0fb0",
				"-node-name",
				"${node.unique.name}",
			},
			"network_mode": "host",
			"cap_add": []string{
				"NET_ADMIN",
				"NET_RAW",
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(40),
		nomadic.WithLoggingTag(a.beaconName))

	return clients.DeployJobs(ctx, beaconJob)
}

func (a *App) installConsulConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	httpName := "detect-presence"
	httpDefaults := nomadic.NewServiceDefaults(httpName, "http")

	if _, _, err := clients.Consul.ConfigEntries().Set(httpDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", httpName, err)
	}

	httpIntentions := nomadic.NewServiceIntentions(httpName,
		nomadic.AppAwareIntention("ingress-http",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/app/")),
			nomadic.DenyAllHTTP()),
		nomadic.DenyIntention("*"))

	if _, _, err := clients.Consul.ConfigEntries().Set(httpIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", httpName, err)
	}

	grpcName := httpName + "-grpc"
	grpcDefaults := nomadic.NewServiceDefaults(grpcName, "grpc")
	if _, _, err := clients.Consul.ConfigEntries().Set(grpcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", grpcName, err)
	}

	grpcIntentions := nomadic.NewServiceIntentions(grpcName,
		nomadic.AppAwareIntention("homebase-bot",
			nomadic.DenyHTTP(nomadic.HTTPPathExact("/TripsService/RecordTrips")),
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/TripsService/"))),
		nomadic.AppAwareIntention("homebase-api",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/TripsService/")),
			nomadic.DenyAllHTTP()),
		nomadic.DenyIntention("*"))

	if _, _, err := clients.Consul.ConfigEntries().Set(grpcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", grpcName, err)
	}

	return nil
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	httpName := "detect-presence"
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

	if _, _, err := clients.Nomad.Jobs().Deregister(a.beaconName, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.beaconName, err)
	}

	return nil
}

//go:embed presence.hcl
var vaultPolicy string
