package ingress

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed load-balancer.conf
var nginxConf string
var nginxTemplate *template.Template

func init() {
	nginxTemplate = template.Must(
		template.New("load_balancer").
			Delims("<<", ">>").
			Parse(nginxConf))
}

type nginxTemplateInput struct {
	VirtualHosts []virtualHost
	Upstreams    []upstream
}

func nginxConfig() (string, error) {
	input := nginxTemplateInput{
		VirtualHosts: virtualHosts,
		Upstreams:    sortedUpstreams(),
	}

	var nginxResult strings.Builder
	if err := nginxTemplate.Execute(&nginxResult, input); err != nil {
		panic(err)
	}

	return nginxResult.String(), nil
}
