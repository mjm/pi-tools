package blocky

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const imageRepo = "spx01/blocky"
const imageVersion = "sha256:59b3661951c28db0eecd9bb2e673c798d7c861d286e7713665da862e5254c477"

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
	job := &nomadapi.Job{
		ID:          &a.name,
		Datacenters: nomadic.DefaultDatacenters,
		Priority:    nomadic.Int(90),
		TaskGroups: []*nomadapi.TaskGroup{
			{
				Name:  nomadic.String("blocky"),
				Count: nomadic.Int(1),
				Networks: []*nomadapi.NetworkResource{
					{
						DNS: nomadic.DefaultDNS,
						DynamicPorts: []nomadapi.Port{
							{
								Label: "dns",
							},
							{
								Label: "http",
							},
						},
					},
				},
				Services: []*nomadapi.Service{
					{
						Name:      a.name,
						PortLabel: "dns",
						TaskName:  "blocky",
						Tags:      []string{"dns"},
						Checks: []nomadapi.ServiceCheck{
							{
								Type:    "script",
								Command: "dig",
								Args: []string{
									"@${NOMAD_IP_dns}",
									"-p",
									"${NOMAD_HOST_PORT_dns}",
									"google.com",
								},
								Interval: 30 * time.Second,
								Timeout:  5 * time.Second,
							},
						},
					},
					{
						Name:      "blocky",
						PortLabel: "http",
						Tags:      []string{"http"},
						Meta: map[string]string{
							"metrics_path": "/metrics",
						},
						Checks: []nomadapi.ServiceCheck{
							{
								Type:                 "http",
								Path:                 "/",
								Interval:             30 * time.Second,
								Timeout:              5 * time.Second,
								SuccessBeforePassing: 3,
							},
						},
					},
				},
				Tasks: []*nomadapi.Task{
					{
						Name:   "blocky",
						Driver: "docker",
						Config: map[string]interface{}{
							"image": imageRepo + "@" + imageVersion,
							"args": []string{
								"/app/blocky",
								"--config",
								"${NOMAD_TASK_DIR}/config.yaml",
							},
							"logging": nomadic.Logging("blocky"),
							"ports":   []string{"dns", "http"},
						},
						Resources: &nomadapi.Resources{
							CPU:      nomadic.Int(200),
							MemoryMB: nomadic.Int(75),
						},
						Templates: []*nomadapi.Template{
							{
								EmbeddedTmpl: nomadic.String(configFile),
								DestPath:     nomadic.String("local/config.yaml"),
								ChangeMode:   nomadic.String("restart"),
							},
						},
					},
				},
			},
		},
	}
	resp, _, err := clients.Nomad.Jobs().Plan(job, true, nil)
	if err != nil {
		return fmt.Errorf("planning %s job: %w", *job.ID, err)
	}
	if resp.Diff.Type == "None" {
		return nil
	}

	if _, _, err := clients.Nomad.Jobs().Register(job, nil); err != nil {
		return fmt.Errorf("registering %s job: %w", *job.ID, err)
	}
	return nil
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return err
	}
	return nil
}

var _ nomadic.Deployable = (*App)(nil)

//go:embed config.yaml
var configFile string
