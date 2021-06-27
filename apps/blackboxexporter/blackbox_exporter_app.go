package blackboxexporter

import (
	"context"
	_ "embed"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo    = "prom/blackbox-exporter"
	imageVersion = "sha256:7c3e8d34768f2db17dce800b0b602196871928977f205bbb8ab44e95a8821be5"
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
	job := nomadic.NewJob(a.name, 70)
	tg := nomadic.AddTaskGroup(job, "blackbox-exporter", 1)

	nomadic.AddPort(tg, nomadapi.Port{Label: "http"})
	nomadic.AddService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "http",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "blackbox-exporter",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"args": []string{
				"--config.file=${NOMAD_TASK_DIR}/blackbox.yml",
				"--web.listen-address=:${NOMAD_PORT_http}",
			},
			"ports": []string{
				"http",
			},
			"network_mode": "host",
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: &configFile,
				DestPath:     nomadic.String("local/blackbox.yml"),
			},
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	panic("implement me")
}

//go:embed blackbox.yml
var configFile string
