package otel

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
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

func (a *App) collectorName() string {
	return a.name + "-collector"
}

func (a *App) Install(ctx context.Context, clients nomadic.Clients) error {
	collectorName := a.collectorName()

	if err := clients.Vault.Sys().PutPolicy(collectorName, vaultPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", collectorName, err)
	}

	job := nomadic.NewSystemJob(collectorName, 70)
	tg := nomadic.AddTaskGroup(job, "otel-collector", 1)
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "healthcheck",
		To:    13133,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "jaeger_thrift",
		To:    14268,
		Value: 14268,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "zipkin",
		To:    9411,
		Value: 9411,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "otlp_grpc",
		To:    4317,
		Value: 4317,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "otlp_grpc_2",
		To:    55680,
		Value: 55680,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "otlp_http",
		To:    55681,
		Value: 55681,
	})
	nomadic.AddPort(tg, nomadapi.Port{
		Label: "metrics",
		To:    8888,
	})

	nomadic.AddService(tg, &nomadapi.Service{
		Name:      collectorName,
		PortLabel: "otlp_grpc",
		Tags:      []string{"grpc"},
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/",
				PortLabel:            "healthcheck",
				Interval:             15 * time.Second,
				Timeout:              3 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	},
		nomadic.WithMetricsScraping("/metrics"),
		nomadic.WithMetricsPort("metrics"))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "otel-collector",
		Config: map[string]interface{}{
			"image":   nomadic.Image("mmoriarity/opentelemetry-collector", "latest"),
			"command": "/otelcol",
			"args": []string{
				"--config",
				"${NOMAD_SECRETS_DIR}/config.yaml",
				"--mem-ballast-size-mib",
				"50",
			},
			"ports": []string{
				"jaeger_thrift",
				"zipkin",
				"otlp_grpc",
				"otlp_grpc_2",
				"otlp_http",
				"healthcheck",
				"metrics",
			},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: &configFile,
				DestPath:     nomadic.String("secrets/config.yaml"),
				ChangeMode:   nomadic.String("restart"),
			},
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(150),
		nomadic.WithLoggingTag(collectorName),
		nomadic.WithVaultPolicies(collectorName),
		nomadic.WithVaultChangeNoop())

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.collectorName()); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.collectorName(), err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.collectorName(), false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.collectorName(), err)
	}

	return nil
}

//go:embed otel-collector.hcl
var vaultPolicy string

//go:embed otel-collector-config.yaml
var configFile string
