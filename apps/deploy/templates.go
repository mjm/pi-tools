package deploy

import (
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

var taskTemplates = []*nomadapi.Template{
	{
		EmbeddedTmpl: nomadic.String(`{{ with secret "kv/deploy" }}{{ .Data.data.github_token }}{{ end }}`),
		DestPath:     nomadic.String("secrets/github-token"),
		ChangeMode:   nomadic.String("restart"),
	},
	{
		EmbeddedTmpl: nomadic.String(`
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ end }}
`),
		DestPath: nomadic.String("secrets/nomad.crt"),
	},
	{
		EmbeddedTmpl: nomadic.String(`
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.private_key }}
{{ end }}
`),
		DestPath: nomadic.String("secrets/nomad.key"),
	},
	{
		EmbeddedTmpl: nomadic.String(`
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
`),
		DestPath: nomadic.String("secrets/nomad.ca.crt"),
	},
	{
		EmbeddedTmpl: nomadic.String(`
{{ with secret "nomad/creds/deploy" }}
NOMAD_TOKEN={{ .Data.secret_id }}
{{ end }}
{{ with secret "kv/pushover" }}
PUSHOVER_USER_KEY={{ .Data.data.user_key }}
PUSHOVER_TOKEN={{ .Data.data.token }}
{{ end }}
{{ with secret "kv/deploy" }}
AWS_ACCESS_KEY_ID=deploy
AWS_SECRET_ACCESS_KEY={{ .Data.data.minio_secret_key }}
{{ end }}
`),
		DestPath:   nomadic.String("secrets/deploy.env"),
		Envvars:    nomadic.Bool(true),
		ChangeMode: nomadic.String("restart"),
	},
}
