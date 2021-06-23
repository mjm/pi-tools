package grafana

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
	imageRepo    = "grafana/grafana"
	imageVersion = "sha256:f0817ecbf8dcf33e10cca2245bd25439433c441189bbe1ce935ac61d05f9cc6f"
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
	if err := clients.Vault.Sys().PutPolicy("grafana", vaultPolicy); err != nil {
		return fmt.Errorf("updating grafana vault policy: %w", err)
	}

	svcDefaults := &consulapi.ServiceConfigEntry{
		Kind:     consulapi.ServiceDefaults,
		Name:     a.name,
		Protocol: "http",
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", a.name, err)
	}

	svcIntentions := &consulapi.ServiceIntentionsConfigEntry{
		Kind: consulapi.ServiceIntentions,
		Name: a.name,
		Sources: []*consulapi.SourceIntention{
			{
				Name:       "ingress-http",
				Precedence: 9,
				Type:       consulapi.IntentionSourceConsul,
				Permissions: []*consulapi.IntentionPermission{
					{
						Action: consulapi.IntentionActionAllow,
						HTTP: &consulapi.IntentionHTTPPermission{
							PathPrefix: "/",
						},
					},
				},
			},
			{
				Action:     consulapi.IntentionActionDeny,
				Name:       "*",
				Precedence: 8,
				Type:       consulapi.IntentionSourceConsul,
			},
		},
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", err)
	}

	job := &nomadapi.Job{
		ID:          &a.name,
		Datacenters: nomadic.DefaultDatacenters,
		Priority:    nomadic.Int(70),
		TaskGroups: []*nomadapi.TaskGroup{
			{
				Name:  nomadic.String("grafana"),
				Count: nomadic.Int(3),
				Networks: []*nomadapi.NetworkResource{
					{
						Mode: "bridge",
						DNS:  nomadic.DefaultDNS,
						DynamicPorts: []nomadapi.Port{
							{
								Label: "health",
							},
							{
								Label: "expose",
							},
							{
								Label: "envoy_metrics",
								To:    9102,
							},
						},
					},
				},
				Services: []*nomadapi.Service{
					{
						Name:      a.name,
						PortLabel: "3000",
						Meta: map[string]string{
							"metrics_path":       "/metrics",
							"envoy_metrics_port": "${NOMAD_HOST_PORT_envoy_metrics}",
							"metrics_port":       "${NOMAD_HOST_PORT_expose}",
						},
						Checks: []nomadapi.ServiceCheck{
							{
								Type:      "http",
								Path:      "/api/health",
								Interval:  15 * time.Second,
								Timeout:   3 * time.Second,
								Expose:    true,
								PortLabel: "health",
							},
						},
						Connect: &nomadapi.ConsulConnect{
							SidecarService: &nomadapi.ConsulSidecarService{
								Proxy: &nomadapi.ConsulProxy{
									ExposeConfig: &nomadapi.ConsulExposeConfig{
										Path: []*nomadapi.ConsulExposePath{
											{
												Path:          "/metrics",
												Protocol:      "http",
												LocalPathPort: 3000,
												ListenerPort:  "expose",
											},
										},
									},
									Upstreams: []*nomadapi.ConsulUpstream{
										nomadic.ConsulUpstream("loki", 3100),
									},
									Config: map[string]interface{}{
										"envoy_prometheus_bind_addr": "0.0.0.0:9102",
									},
								},
							},
						},
					},
				},
				Tasks: []*nomadapi.Task{
					{
						Name:   "grafana",
						Driver: "docker",
						Config: map[string]interface{}{
							"image":   nomadic.Image(imageRepo, imageVersion),
							"logging": nomadic.Logging("grafana"),
						},
						Resources: &nomadapi.Resources{
							CPU:      nomadic.Int(200),
							MemoryMB: nomadic.Int(100),
						},
						Env: map[string]string{
							"GF_PATHS_CONFIG":       "${NOMAD_SECRETS_DIR}/grafana.ini",
							"GF_PATHS_PROVISIONING": "${NOMAD_TASK_DIR}/provisioning",
						},
						Vault: &nomadapi.Vault{
							Policies: []string{"grafana"},
						},
						Templates: taskTemplates(),
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
	if err := clients.Vault.Sys().DeletePolicy("grafana"); err != nil {
		return fmt.Errorf("deleting grafana vault policy: %w", err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s job: %w", a.name, err)
	}

	return nil
}

var _ nomadic.Deployable = (*App)(nil)
