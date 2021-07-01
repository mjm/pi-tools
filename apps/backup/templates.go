package backup

import (
	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

var tarsnapKeyTemplate = &nomadapi.Template{
	EmbeddedTmpl: nomadic.String(`{{ with secret "kv/tarsnap" }}{{ .Data.data.key | base64Decode }}{{ end }}`),
	DestPath:     nomadic.String("secrets/tarsnap.key"),
}

var borgSSHKeyTemplate = &nomadapi.Template{
	EmbeddedTmpl: nomadic.String(`{{ with secret "kv/borg" }}{{ .Data.data.private_key }}{{ end }}
`),
	DestPath: nomadic.String("secrets/id_rsa"),
	Perms:    nomadic.String("0600"),
}
