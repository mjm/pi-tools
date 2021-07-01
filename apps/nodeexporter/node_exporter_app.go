package nodeexporter

import (
	"context"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "prom/node-exporter"
	imageVersion = "sha256:eb80355f0ff0a0a0f0342303cd694af28e2820d688f416049d7be7d1760a0b33"
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
	job := nomadic.NewSystemJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "node-exporter", 1)

	nomadic.AddPort(tg, nomadapi.Port{Label: "http"})

	nomadic.AddService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "http",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/",
				Interval:             30 * time.Second,
				Timeout:              5 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "node-exporter",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"--web.listen-address=:${NOMAD_PORT_http}",
				"--path.procfs=/host/proc",
				"--path.sysfs=/host/sys",
				"--path.rootfs=/host/root",
				"--collector.processes",
				"--collector.systemd",
				"--collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+|var/lib/nomad/.+|run/docker/.+|snap/.+)($|/)",
				"--collector.netclass.ignored-devices=^veth",
				"--collector.netdev.device-exclude=^veth",
			},
			"privileged":   true,
			"pid_mode":     "host",
			"network_mode": "host",
			"mount": []map[string]interface{}{
				{
					"type":   "bind",
					"target": "/host/root",
					"source": "/",
					"bind_options": map[string]interface{}{
						"propagation": "rslave", // :(
					},
				},
				{
					"type":   "bind",
					"target": "/host/proc",
					"source": "/proc",
					"bind_options": map[string]interface{}{
						"propagation": "rslave", // :(
					},
				},
				{
					"type":   "bind",
					"target": "/host/sys",
					"source": "/sys",
					"bind_options": map[string]interface{}{
						"propagation": "rslave", // :(
					},
				},
				{
					"type":   "bind",
					"target": "/run/systemd",
					"source": "/run/systemd",
					"bind_options": map[string]interface{}{
						"propagation": "rslave", // :(
					},
				},
				{
					"type":   "bind",
					"target": "/var/run/dbus",
					"source": "/var/run/dbus",
				},
			},
		},
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}
