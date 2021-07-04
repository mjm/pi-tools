package nut

import (
	"context"
	"fmt"
	"strings"

	nomad "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo = "instantlinux/nut-upsd"
	// nut 2.7.4-r8
	imageVersion = "sha256:f4948378c7e7c27c857e822aac174854389b8607fbe937538ef647d9cade69b0"

	exporterImageRepo = "mmoriarity/nut_exporter"
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
	tg := nomadic.AddTaskGroup(job, "nut", 1)

	// this needs to connect to the Tripplite UPS over USB,
	// and it's only plugged in to this machine.
	job.Constrain(&nomad.Constraint{
		LTarget: "${node.unique.name}",
		Operand: "=",
		RTarget: "raspberrypi",
	})

	nomadic.AddPort(tg, nomad.Port{
		Label: "upsd",
		Value: 3493,
		To:    3493,
	})

	nomadic.AddPort(tg, nomad.Port{
		Label: "metrics",
		To:    9199,
	})

	nomadic.AddService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "upsd",
		Meta: map[string]string{
			"metrics_port": "${NOMAD_HOST_PORT_metrics}",
		},
	}, nomadic.WithMetricsScraping("/ups_metrics"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "upsd",
		Config: map[string]interface{}{
			"image":      nomadic.Image(imageRepo, imageVersion),
			"ports":      []string{"upsd"},
			"privileged": true,
			"volumes": []string{
				"secrets:/run/secrets",
			},
		},
		Env: map[string]string{
			"VENDORID": "09ae",
			"SERIAL":   "3031BV4OM882401020",
			// terrible hack to get the productid into the config file
			"POLLINTERVAL": "2\n        productid = 3024",
		},
		Templates: []*nomad.Template{
			{
				// TODO store a real random password in Vault
				EmbeddedTmpl: nomadic.String(`asdf`),
				DestPath:     nomadic.String("secrets/nut-upsd-password"),
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(a.name))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "nut-exporter",
		Config: map[string]interface{}{
			"image":   nomadic.Image(exporterImageRepo, "latest"),
			"command": "/nut_exporter",
			"args": []string{
				"--nut.server=${NOMAD_IP_upsd}",
				"--nut.vars_enable=" + strings.Join(enabledVariables, ","),
			},
			"ports": []string{"metrics"},
		},
	},
		nomadic.WithCPU(50),
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
