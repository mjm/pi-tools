package presence

import (
	"context"
	"fmt"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const imageRepo = "mmoriarity/beacon-srv"

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
	beaconJob := nomadic.NewSystemJob(a.beaconName, 50)
	tg := nomadic.AddTaskGroup(beaconJob, "beacon", 1)

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "beacon-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(imageRepo, "latest"),
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

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.beaconName, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.beaconName, err)
	}

	return nil
}
