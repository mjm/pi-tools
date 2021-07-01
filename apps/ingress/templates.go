package ingress

import (
	_ "embed"
	"strings"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

//go:embed ingress.hcl
var vaultPolicy string

var extraCertNames = []string{
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
