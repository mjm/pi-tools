package backup

import (
	_ "embed"

	"github.com/hashicorp/nomad/api"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

func (a *App) createBorgJob() *api.Job {
	job := nomadic.NewBatchJob(a.borgName(), 70)
	job.AddPeriodicConfig(&api.PeriodicConfig{
		Spec:            nomadic.String("30 */4 * * *"),
		ProhibitOverlap: nomadic.Bool(true),
	})
	tg := nomadic.AddTaskGroup(job, "backup", 1)

	a.addCommonTasks(job, tg)

	nomadic.AddTask(tg, &api.Task{
		Name: "backup",
		Config: map[string]interface{}{
			"image":   nomadic.Image(backupImageRepo, "latest"),
			"command": "/usr/bin/perform-backup",
			"args": []string{
				"-kind",
				"borg",
			},
		},
		Env: map[string]string{
			"BORG_RSH": "ssh -o StrictHostKeyChecking=no -i ${NOMAD_SECRETS_DIR}/id_rsa",
		},
		Templates: []*api.Template{
			{
				EmbeddedTmpl: nomadic.String(`PUSHGATEWAY_URL={{ with service "pushgateway" }}{{ with index . 0 }}http://{{ .Address }}:{{ .Port }}{{ end }}{{ end }}`),
				DestPath:     nomadic.String("local/backup.env"),
				Envvars:      nomadic.Bool(true),
			},
			borgSSHKeyTemplate,
		},
	},
		nomadic.WithCPU(100),
		nomadic.WithMemoryMB(100),
		nomadic.WithLoggingTag(a.borgName()),
		nomadic.WithVaultPolicies(a.borgName()))

	return job
}

func (a *App) borgName() string {
	return a.name + "-borg"
}

//go:embed borg.hcl
var borgPolicy string
