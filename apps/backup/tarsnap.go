package backup

import (
	_ "embed"

	nomad "github.com/hashicorp/nomad/api"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

const (
	tarsnapImageRepo    = "mmoriarity/tarsnap"
	tarsnapImageVersion = "sha256:4deeb35783541c160a09cb7a58489a7bf57bb456f4efab83e0cbd663a60bbf50"
)

func (a *App) createTarsnapJob() *nomad.Job {
	job := nomadic.NewBatchJob(a.tarsnapName(), 70)
	job.AddPeriodicConfig(&nomad.PeriodicConfig{
		Spec:            nomadic.String("0 12 * * *"),
		ProhibitOverlap: nomadic.Bool(true),
	})
	tg := nomadic.AddTaskGroup(job, "backup", 1)

	// only run tarsnap job on this node because it has the tarsnap cache
	// TODO make the tarsnap cache into a host volume and mount it that way
	tg.Constrain(&nomad.Constraint{
		LTarget: "${node.unique.name}",
		Operand: "=",
		RTarget: "raspberrypi",
	})

	a.addCommonTasks(job, tg)

	mysqlDumpArgs := append([]string{
		"mysqldump",
		"--defaults-file=${NOMAD_SECRETS_DIR}/my.cnf",
		"--result-file=${NOMAD_ALLOC_DIR}/data/phabricator.sql",
		"--databases",
	}, phabricatorDatabases...)

	nomadic.AddTask(tg, &nomad.Task{
		Name: "dump-phabricator-dbs",
		Lifecycle: &nomad.TaskLifecycle{
			Hook: nomad.TaskLifecycleHookPrestart,
		},
		Config: map[string]interface{}{
			"image": nomadic.Image(mysqlImageRepo, mysqlImageVersion),
			"args":  mysqlDumpArgs,
		},
		Templates: []*nomad.Template{
			{
				EmbeddedTmpl: nomadic.String(`
[client]
host = mysql.service.consul
{{ with secret "database/creds/phabricator" -}}
user = {{ .Data.username }}
password = {{ .Data.password }}
{{- end }}
`),
				DestPath: nomadic.String("secrets/my.cnf"),
			},
		},
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(500),
		nomadic.WithLoggingTag(a.tarsnapName()),
		nomadic.WithVaultPolicies("phabricator"))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "backup",
		Config: map[string]interface{}{
			"image":   nomadic.Image(backupImageRepo, "latest"),
			"command": "/usr/bin/perform-backup",
			"args": []string{
				"-kind",
				"tarsnap",
			},
			"mount": []map[string]interface{}{
				{
					"type":   "bind",
					"target": "/var/lib/tarsnap/cache",
					"source": "/var/lib/tarsnap/cache",
				},
			},
		},
		Templates: []*nomad.Template{
			{
				EmbeddedTmpl: nomadic.String(`PUSHGATEWAY_URL={{ with service "pushgateway" }}{{ with index . 0 }}http://{{ .Address }}:{{ .Port }}{{ end }}{{ end }}`),
				DestPath:     nomadic.String("local/backup.env"),
				Envvars:      nomadic.Bool(true),
			},
			tarsnapKeyTemplate,
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(100),
		nomadic.WithLoggingTag(a.tarsnapName()),
		nomadic.WithVaultPolicies(a.tarsnapName()))

	nomadic.AddTask(tg, &nomad.Task{
		Name: "prune",
		Lifecycle: &nomad.TaskLifecycle{
			Hook: nomad.TaskLifecycleHookPoststop,
		},
		Config: map[string]interface{}{
			"image":   nomadic.Image(backupImageRepo, "latest"),
			"command": "${NOMAD_TASK_DIR}/prune.sh",
			"mount": []map[string]interface{}{
				{
					"type":   "bind",
					"target": "/var/lib/tarsnap/cache",
					"source": "/var/lib/tarsnap/cache",
				},
			},
		},
		Templates: []*nomad.Template{
			tarsnapKeyTemplate,
			{
				EmbeddedTmpl: &pruneScript,
				DestPath:     nomadic.String("local/prune.sh"),
				Perms:        nomadic.String("0755"),
			},
		},
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(30),
		nomadic.WithLoggingTag(a.tarsnapName()),
		nomadic.WithVaultPolicies(a.tarsnapName()))

	return job
}

func (a *App) createTarsnapDeleteJob() *nomad.Job {
	name := a.tarsnapName() + "-delete"
	job := nomadic.NewBatchJob(name, 70)
	job.ParameterizedJob = &nomad.ParameterizedJobConfig{
		Payload: "required",
	}
	tg := nomadic.AddTaskGroup(job, "tarsnap-delete", 1)

	// only run tarsnap job on this node because it has the tarsnap cache
	// TODO make the tarsnap cache into a host volume and mount it that way
	tg.Constrain(&nomad.Constraint{
		LTarget: "${node.unique.name}",
		Operand: "=",
		RTarget: "raspberrypi",
	})

	nomadic.AddTask(tg, &nomad.Task{
		Name: "delete-backup",
		Config: map[string]interface{}{
			"image":   nomadic.Image(tarsnapImageRepo, tarsnapImageVersion),
			"command": "sh",
			"args": []string{
				"-c",
				"cat ${NOMAD_TASK_DIR}/archives.txt | xargs -n1 tarsnap --keyfile ${NOMAD_SECRETS_DIR}/tarsnap.key --cachedir /var/lib/tarsnap/cache --no-default-config -v -d -f",
			},
			"mount": []map[string]interface{}{
				{
					"type":   "bind",
					"target": "/var/lib/tarsnap/cache",
					"source": "/var/lib/tarsnap/cache",
				},
			},
		},
		DispatchPayload: &nomad.DispatchPayloadConfig{
			File: "archives.txt",
		},
		Templates: []*nomad.Template{
			tarsnapKeyTemplate,
		},
	},
		nomadic.WithCPU(200),
		nomadic.WithMemoryMB(30),
		nomadic.WithLoggingTag(a.tarsnapName()),
		nomadic.WithVaultPolicies(a.tarsnapName()))

	return job
}

func (a *App) tarsnapName() string {
	return a.name + "-tarsnap"
}

//go:embed tarsnap.hcl
var tarsnapPolicy string

//go:embed prune.sh
var pruneScript string
