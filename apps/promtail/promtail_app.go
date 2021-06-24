package promtail

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo = "grafana/promtail"
	// promtail 2.2.1
	imageVersion = "sha256:ca2711bece9b74ce51aad398dedeba706c553f16446a79d0b495573a0060529b"
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
	job := nomadic.NewSystemJob(a.name, 80)
	tg := nomadic.AddTaskGroup(job, "promtail", 1)

	nomadic.AddPort(tg, nomadapi.Port{
		Label: "syslog",
		Value: 3102,
	})

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "3101",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:     "http",
				Path:     "/ready",
				Interval: 15 * time.Second,
				Timeout:  3 * time.Second,
			},
		},
	},
		nomadic.WithMetricsScraping("/metrics"),
		nomadic.WithUpstreams(nomadic.ConsulUpstream("loki", 3100)))

	nomadic.AddService(tg, &nomadapi.Service{
		Name:      "promtail",
		PortLabel: "syslog",
		Tags:      []string{"syslog"},
	})

	tg.Volumes = map[string]*nomadapi.VolumeRequest{
		"run": {
			Name:   "run",
			Type:   "host",
			Source: "promtail_run",
		},
	}

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "promtail",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"-config.file=${NOMAD_TASK_DIR}/promtail.yml",
			},
			"ports": []string{"syslog"},
			"mount": []map[string]interface{}{
				{
					"type":   "bind",
					"source": "/var/log/journal",
					"target": "/var/log/journal",
				},
				{
					"type":   "bind",
					"source": "/run/log/journal",
					"target": "/run/log/journal",
				},
				{
					"source": "/etc/machine-id",
					"target": "/etc/machine-id",
					"type":   "bind",
				},
			},
		},
		VolumeMounts: []*nomadapi.VolumeMount{
			{
				Volume:      nomadic.String("run"),
				Destination: nomadic.String("/run/promtail"),
			},
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: &configFile,
				DestPath:     nomadic.String("local/promtail.yml"),
			},
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name))

	return clients.DeployJobs(ctx, job)
}

//go:embed promtail.yml
var configFile string

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}

var _ nomadic.Deployable = (*App)(nil)
