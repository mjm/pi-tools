package ingress

import (
	"context"
	_ "embed"
	"fmt"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

type App struct {
	name string
}

var _ nomadic.Deployable = (*App)(nil)

func New(name string) *App {
	return &App{
		name: name,
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) Install(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().PutPolicy("ingress", vaultPolicy); err != nil {
		return fmt.Errorf("updating ingress vault policy: %w", err)
	}

	svcDefaults := &consulapi.ServiceConfigEntry{
		Kind:     consulapi.ServiceDefaults,
		Name:     "ingress-http",
		Protocol: "http",
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("updating ingress service defaults: %w", err)
	}

	svcIntentions := &consulapi.ServiceIntentionsConfigEntry{
		Kind: consulapi.ServiceIntentions,
		Name: "ingress-http",
		Sources: []*consulapi.SourceIntention{
			{
				Action:     consulapi.IntentionActionDeny,
				Name:       "*",
				Precedence: 8,
				Type:       consulapi.IntentionSourceConsul,
			},
		},
	}
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("updating ingress service intentions: %w", err)
	}

	job := &nomadapi.Job{
		ID:          &a.name,
		Datacenters: nomadic.DefaultDatacenters,
		Priority:    nomadic.Int(70),
		Update: &nomadapi.UpdateStrategy{
			MaxParallel: nomadic.Int(1),
			Stagger:     nomadic.Duration(30 * time.Second),
		},
		TaskGroups: []*nomadapi.TaskGroup{
			{
				Name:  nomadic.String("ingress"),
				Count: nomadic.Int(2),
				Networks: []*nomadapi.NetworkResource{
					{
						Mode: "bridge",
						DNS:  nomadic.DefaultDNS,
						DynamicPorts: []nomadapi.Port{
							{
								Label:       "health",
								HostNetwork: "default",
							},
						},
						ReservedPorts: []nomadapi.Port{
							{
								Label: "http",
								To:    80,
								Value: 80,
							},
							{
								Label: "https",
								To:    443,
								Value: 443,
							},
						},
					},
				},
				Services: []*nomadapi.Service{
					{
						Name:      "ingress-http",
						PortLabel: "80",
						Checks: []nomadapi.ServiceCheck{
							{
								Type:      "http",
								Path:      "/healthz",
								PortLabel: "health",
								Expose:    true,
								Interval:  15 * time.Second,
								Timeout:   3 * time.Second,
								Header: map[string][]string{
									// Just need to have some Host header
									"Host": {"ingress"},
								},
							},
						},
						Connect: &nomadapi.ConsulConnect{
							SidecarService: &nomadapi.ConsulSidecarService{
								Proxy: &nomadapi.ConsulProxy{
									Upstreams: []*nomadapi.ConsulUpstream{
										nomadic.ConsulUpstream("detect-presence", 2120),
										nomadic.ConsulUpstream("vault-proxy", 2220),
										nomadic.ConsulUpstream("go-links", 4240),
										nomadic.ConsulUpstream("homebase-api", 6460),
										nomadic.ConsulUpstream("homebase", 3001),
										nomadic.ConsulUpstream("grafana", 3000),
									},
								},
							},
						},
					},
					{
						Name:      "ingress-https",
						PortLabel: "443",
					},
				},
				Tasks: []*nomadapi.Task{
					{
						Name:   "nginx",
						Driver: "docker",
						Config: map[string]interface{}{
							"image": "nginx@sha256:763d95e3db66d9bd1bb926c029e5659ee67eb49ff57f83d331de5f5af6d2ae0c",
							"volumes": []string{
								"local:/etc/nginx/conf.d",
								"secrets:/etc/nginx/ssl",
							},
							"logging": nomadic.Logging("ingress"),
						},
						Resources: &nomadapi.Resources{
							CPU:      nomadic.Int(100),
							MemoryMB: nomadic.Int(50),
						},
						Vault: &nomadapi.Vault{
							Policies:   []string{"ingress"},
							ChangeMode: nomadic.String("noop"),
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
	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return err
	}
	return nil
}

//go:embed ingress.hcl
var vaultPolicy string

var extraCertNames = []string{
	"consul",
	"homebase",
}

const certTemplateData = `
{{ with secret "pki-homelab/issue/homelab" "common_name=CERTNAME.homelab" "alt_names=CERTNAME.home.mattmoriarity.com" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
`

func taskTemplates() []*nomadapi.Template {
	nginxResult, err := nginxConfig()
	if err != nil {
		panic(err)
	}

	templates := []*nomadapi.Template{
		{
			EmbeddedTmpl: &nginxResult,
			DestPath:     nomadic.String("local/load-balancer.conf"),
			ChangeMode:   nomadic.String("signal"),
			ChangeSignal: nomadic.String("SIGHUP"),
		},
		{
			EmbeddedTmpl: nomadic.String(`
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
`),
			DestPath:     nomadic.String("secrets/nomad.pem"),
			ChangeMode:   nomadic.String("signal"),
			ChangeSignal: nomadic.String("SIGHUP"),
		},
		{
			EmbeddedTmpl: nomadic.String(`
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
`),
			DestPath:     nomadic.String("secrets/nomad.ca.crt"),
			ChangeMode:   nomadic.String("signal"),
			ChangeSignal: nomadic.String("SIGHUP"),
		},
	}

	certNames := append([]string{}, extraCertNames...)
	for _, vhost := range virtualHosts {
		certNames = append(certNames, vhost.Name)
	}

	for _, certName := range certNames {
		data := strings.ReplaceAll(certTemplateData, "CERTNAME", certName)
		templates = append(templates, &nomadapi.Template{
			EmbeddedTmpl: &data,
			DestPath:     nomadic.String("secrets/" + certName + ".homelab.pem"),
			ChangeMode:   nomadic.String("signal"),
			ChangeSignal: nomadic.String("SIGHUP"),
		})
	}

	return templates
}
