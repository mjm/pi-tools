package adminer

import (
	"context"
	_ "embed"
	"fmt"

	consul "github.com/hashicorp/consul/api"
	nomad "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	imageRepo = "adminer"
	// adminer 4.8.1
	imageVersion = "sha256:c98cafe3e363b3a9008361515fd7a7e714840d69577d105967499fb0ae7475f8"
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

	svcDefaults := nomadic.NewServiceDefaults(a.name, "http")
	if _, _, err := clients.Consul.ConfigEntries().Set(svcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", svcDefaults.Name, err)
	}

	svcIntentions := nomadic.NewServiceIntentions(a.name,
		nomadic.AppAwareIntention("ingress-http",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/"))),
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(svcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", svcIntentions.Name, err)
	}

	job := nomadic.NewJob(a.name, 50)
	// It's important that this only has 1 task, because it stores session data in files, so if there's multiple,
	// then we load-balance between them and don't have consistent session data. We could probably also just IP hash
	// in the ingress, but there's not much reason to need redundancy for something like this.
	tg := nomadic.AddTaskGroup(job, "adminer", 1)

	nomadic.AddConnectService(tg, &nomad.Service{
		Name:      a.name,
		PortLabel: "8080",
	})

	nomadic.AddTask(tg, &nomad.Task{
		Name: "adminer",
		Config: map[string]interface{}{
			"image": nomadic.Image(imageRepo, imageVersion),
			"volumes": []string{
				"secrets/plugins:/var/www/html/plugins-enabled",
			},
		},
		Templates: []*nomad.Template{
			{
				EmbeddedTmpl: &loginStaticPlugin,
				DestPath:     nomadic.String("secrets/plugins/login-static.php"),
			},
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(100),
		nomadic.WithLoggingTag(a.name),
		nomadic.WithVaultPolicies(a.name))

	return clients.DeployJobs(ctx, job)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consul.ServiceDefaults, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", a.name, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consul.ServiceIntentions, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", a.name, err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.name, false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.name, err)
	}

	return nil
}

//go:embed adminer.hcl
var vaultPolicy string

//go:embed login-static.php
var loginStaticPlugin string
