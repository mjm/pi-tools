package presence

import (
	"context"
	_ "embed"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
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

	beaconJob := nomadic.NewSystemJob(a.beaconName, 50)
	tg := nomadic.AddTaskGroup(beaconJob, "beacon", 1)

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
