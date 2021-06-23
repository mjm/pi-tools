package grafana

import (
	_ "embed"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/apps/grafana/dashboards"
	"github.com/mjm/pi-tools/pkg/nomadic"
)

var (
	//go:embed grafana.hcl
	vaultPolicy string

	//go:embed grafana.ini
	grafanaConfig string

	//go:embed datasources.yaml
	datasourcesConfig string

	//go:embed dashboards.yaml
	dashboardsConfig string
)

func taskTemplates() []*nomadapi.Template {
	templates := []*nomadapi.Template{
		{
			EmbeddedTmpl: &grafanaConfig,
			DestPath:     nomadic.String("secrets/grafana.ini"),
			ChangeMode:   nomadic.String("restart"),
		},
		{
			EmbeddedTmpl: &datasourcesConfig,
			DestPath:     nomadic.String("local/provisioning/datasources/datasources.yaml"),
			ChangeMode:   nomadic.String("restart"),
		},
		{
			EmbeddedTmpl: &dashboardsConfig,
			DestPath:     nomadic.String("local/provisioning/dashboards/dashboards.yaml"),
			ChangeMode:   nomadic.String("restart"),
		},
	}

	dashes, err := dashboards.All.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, entry := range dashes {
		contents, err := dashboards.All.ReadFile(entry.Name())
		if err != nil {
			panic(err)
		}

		templates = append(templates, &nomadapi.Template{
			EmbeddedTmpl: nomadic.String(string(contents)),
			DestPath:     nomadic.String("local/dashboards/" + entry.Name()),
			// prevent interpreting blocks delimited by '{{' and '}}' as consul templates
			LeftDelim: nomadic.String("do_not_substitute"),
		})
	}

	return templates
}
