package deploy

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
	imageRepo = "mmoriarity/deploy-srv"
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

	grpcName := a.name + "-grpc"
	svcDefaults := nomadic.NewServiceDefaults(grpcName, "grpc")
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", grpcName, err)
	}

	svcIntentions := nomadic.NewServiceIntentions(grpcName,
		nomadic.AppAwareIntention("homebase-api",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/DeployService/")),
			nomadic.DenyAllHTTP()),
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", grpcName, err)
	}

	job := nomadic.NewJob(a.name, 60)
	tg := nomadic.AddTaskGroup(job, "deploy", 2)
	tg.Update = &nomadapi.UpdateStrategy{
		MaxParallel: nomadic.Int(1),
	}

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "8480",
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
		PortLabel: "8481",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "deploy-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
			"command": "/deploy-srv",
			"args": []string{
				"-leader-elect",
				"-minio-url",
				"http://minio.service.consul:9000",
			},
		},
		// If deploy-srv is updating itself, we want to keep the old version alive until it has finished
		// the in-progress deploy. Otherwise, we leave around in-progress deploys in GitHub and don't get
		// the right report.
		KillTimeout: nomadic.Duration(10 * time.Minute),
		Env: map[string]string{
			"CONSUL_HTTP_ADDR":  "${attr.unique.network.ip-address}:8500",
			"NOMAD_ADDR":        "https://nomad.service.consul:4646",
			"NOMAD_CACERT":      "${NOMAD_SECRETS_DIR}/nomad.ca.crt",
			"NOMAD_CLIENT_CERT": "${NOMAD_SECRETS_DIR}/nomad.crt",
			"NOMAD_CLIENT_KEY":  "${NOMAD_SECRETS_DIR}/nomad.key",
			"VAULT_ADDR":        "http://vault.service.consul:8200",
		},
		Templates: taskTemplates,
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(100),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name),
		nomadic.WithTracingEnv())

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name, err)
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

//go:embed deploy.hcl
var vaultPolicy string
