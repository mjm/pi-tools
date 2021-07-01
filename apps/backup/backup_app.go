package backup

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
	consulImageRepo    = "consul"
	consulImageVersion = "sha256:7b878010be55876f2dd419e0e95aad54cd87ae078d5de54e232e4135eb1069c6"

	postgresImageRepo    = "postgres"
	postgresImageVersion = "sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"

	mysqlImageRepo    = "mysql/mysql-server"
	mysqlImageVersion = "sha256:b33c6e23c8678e29a43ae7cad47cd6bbead6e39c911c5a7b2b6d943cd42b2944"

	backupImageRepo = "mmoriarity/perform-backup"

	serviceImageRepo = "mmoriarity/backup-srv"
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
	if err := clients.Vault.Sys().PutPolicy(a.name, backupPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", a.name, err)
	}

	if err := clients.Vault.Sys().PutPolicy(a.borgName(), borgPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", a.borgName(), err)
	}

	if err := clients.Vault.Sys().PutPolicy(a.tarsnapName(), tarsnapPolicy); err != nil {
		return fmt.Errorf("updating %s vault policy: %w", a.tarsnapName(), err)
	}

	if err := a.installConfigEntries(ctx, clients); err != nil {
		return err
	}

	borgJob := a.createBorgJob()
	tarsnapJob := a.createTarsnapJob()
	tarsnapDeleteJob := a.createTarsnapDeleteJob()
	serviceJob := a.createServiceJob()

	return clients.DeployJobs(ctx, borgJob, tarsnapJob, tarsnapDeleteJob, serviceJob)
}

func (a *App) Uninstall(ctx context.Context, clients nomadic.Clients) error {
	if err := clients.Vault.Sys().DeletePolicy(a.name); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.name, err)
	}

	if err := clients.Vault.Sys().DeletePolicy(a.borgName()); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.borgName(), err)
	}

	if err := clients.Vault.Sys().DeletePolicy(a.tarsnapName()); err != nil {
		return fmt.Errorf("deleting %s vault policy: %w", a.tarsnapName(), err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, a.name, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", a.name, err)
	}

	grpcName := a.name + "-grpc"
	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceDefaults, grpcName, nil); err != nil {
		return fmt.Errorf("deleting %s service defaults: %w", grpcName, err)
	}

	if _, err := clients.Consul.ConfigEntries().Delete(consulapi.ServiceIntentions, grpcName, nil); err != nil {
		return fmt.Errorf("deleting %s service intentions: %w", grpcName, err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.borgName(), false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.borgName(), err)
	}

	if _, _, err := clients.Nomad.Jobs().Deregister(a.tarsnapName(), false, nil); err != nil {
		return fmt.Errorf("deregistering %s nomad job: %w", a.tarsnapName(), err)
	}

	return nil
}

func (a *App) installConfigEntries(ctx context.Context, clients nomadic.Clients) error {
	httpName := a.name
	httpDefaults := nomadic.NewServiceDefaults(httpName, "http")
	if _, _, err := clients.Consul.ConfigEntries().Set(httpDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", httpName, err)
	}

	grpcName := a.name + "-grpc"
	grpcDefaults := nomadic.NewServiceDefaults(grpcName, "grpc")
	if _, _, err := clients.Consul.ConfigEntries().Set(grpcDefaults, nil); err != nil {
		return fmt.Errorf("setting %s service defaults: %w", grpcName, err)
	}

	grpcIntentions := nomadic.NewServiceIntentions(grpcName,
		nomadic.AppAwareIntention("homebase-api",
			nomadic.AllowHTTP(nomadic.HTTPPathPrefix("/BackupService/")),
			nomadic.DenyAllHTTP()),
		nomadic.DenyIntention("*"))
	if _, _, err := clients.Consul.ConfigEntries().Set(grpcIntentions, nil); err != nil {
		return fmt.Errorf("setting %s service intentions: %w", grpcName, err)
	}

	return nil
}

func (a *App) createServiceJob() *nomadapi.Job {
	name := a.name + "-srv"
	job := nomadic.NewJob(name, 50)
	tg := nomadic.AddTaskGroup(job, "backup", 2)
	tg.Update = &nomadapi.UpdateStrategy{
		MaxParallel: nomadic.Int(1),
	}

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name,
		PortLabel: "2320",
		Checks: []nomadapi.ServiceCheck{
			{
				Type:                 "http",
				Path:                 "/healthz",
				Interval:             15 * time.Second,
				Timeout:              3 * time.Second,
				SuccessBeforePassing: 3,
			},
		},
	}, nomadic.WithMetricsScraping("/metrics"))

	nomadic.AddConnectService(tg, &nomadapi.Service{
		Name:      a.name + "-grpc",
		PortLabel: "2321",
	})

	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "backup-srv",
		Config: map[string]interface{}{
			"image":   nomadic.Image(serviceImageRepo, "latest"),
			"command": "/backup-srv",
			"args": []string{
				"-tarsnap-keyfile",
				"${NOMAD_SECRETS_DIR}/tarsnap.key",
			},
		},
		Env: map[string]string{
			"BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK": "yes",

			"BORG_RSH": "ssh -o StrictHostKeyChecking=no -i ${NOMAD_SECRETS_DIR}/id_rsa",
		},
		Templates: []*nomadapi.Template{
			tarsnapKeyTemplate,
			borgSSHKeyTemplate,
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(100),
		nomadic.WithLoggingTag(name),
		nomadic.WithVaultPolicies(a.borgName(), a.tarsnapName()),
		nomadic.WithTracingEnv())

	return job
}

func (a *App) addCommonTasks(job *nomadapi.Job, tg *nomadapi.TaskGroup) {
	nomadic.AddTask(tg, &nomadapi.Task{
		Name: "consul-snapshot",
		Lifecycle: &nomadapi.TaskLifecycle{
			Hook: nomadapi.TaskLifecycleHookPrestart,
		},
		Config: map[string]interface{}{
			"image":   nomadic.Image(consulImageRepo, consulImageVersion),
			"command": "/bin/sh",
			"args": []string{
				"-c",
				"consul snapshot save ${NOMAD_ALLOC_DIR}/data/consul.snap",
			},
			"network_mode": "host",
		},
		Templates: []*nomadapi.Template{
			{
				EmbeddedTmpl: nomadic.String(`CONSUL_HTTP_TOKEN={{ with secret "consul/creds/backup" }}{{ .Data.token }}{{ end }}`),
				DestPath:     nomadic.String("secrets/consul.env"),
				Envvars:      nomadic.Bool(true),
			},
		},
	},
		nomadic.WithCPU(50),
		nomadic.WithMemoryMB(50),
		nomadic.WithLoggingTag(*job.ID),
		nomadic.WithVaultPolicies(a.name))

	for _, db := range pgDatabases {
		nomadic.AddTask(tg, &nomadapi.Task{
			Name: fmt.Sprintf("dump-%s-db", db.Name),
			Lifecycle: &nomadapi.TaskLifecycle{
				Hook: nomadapi.TaskLifecycleHookPrestart,
			},
			Config: map[string]interface{}{
				"image":   nomadic.Image(postgresImageRepo, postgresImageVersion),
				"command": "pg_dump",
				"args": []string{
					"--host=10.0.2.102",
					"--dbname=" + db.DBName(),
					"--file=${NOMAD_ALLOC_DIR}/data/" + db.DBName() + ".sql",
				},
				"network_mode": "host", // is this actually needed?
			},
			Templates: []*nomadapi.Template{
				{
					EmbeddedTmpl: nomadic.String(`
{{ with secret "database/creds/` + db.VaultPolicy() + `" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
`),
					DestPath: nomadic.String("secrets/db.env"),
					Envvars:  nomadic.Bool(true),
				},
			},
		},
			nomadic.WithCPU(50),
			nomadic.WithMemoryMB(50),
			nomadic.WithLoggingTag(*job.ID),
			nomadic.WithVaultPolicies(db.VaultPolicy()))
	}
}

//go:embed backup.hcl
var backupPolicy string
