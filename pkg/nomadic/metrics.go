package nomadic

import (
	"fmt"
	"strconv"

	nomadapi "github.com/hashicorp/nomad/api"
)

type MetricsScrapeConfig struct {
	Path             string
	PortLabel        string
	EnvoyMetricsPort int
	OnlyEnvoyMetrics bool
}

func EnableMetricsScraping(tg *nomadapi.TaskGroup, svc *nomadapi.Service, cfg *MetricsScrapeConfig) error {
	if cfg == nil {
		cfg = &MetricsScrapeConfig{}
	}

	i := getServiceIndex(tg, svc)
	if i < 0 {
		return fmt.Errorf("service %s is not part of task group %s", svc.Name, *tg.Name)
	}

	hasConnectSidecar := svc.Connect != nil && svc.Connect.SidecarService != nil

	if svc.Meta == nil {
		svc.Meta = map[string]string{}
	}
	if !cfg.OnlyEnvoyMetrics {
		if cfg.Path == "" {
			cfg.Path = "/metrics"
		}
		svc.Meta["metrics_path"] = cfg.Path

		if hasConnectSidecar && cfg.PortLabel == "" {
			// connect services must use some port to expose their metrics through envoy
			cfg.PortLabel = "expose"
		}
		if cfg.PortLabel != "" {
			svc.Meta["metrics_port"] = fmt.Sprintf("${NOMAD_HOST_PORT_%s}", cfg.PortLabel)
		}
	}

	if !hasConnectSidecar {
		// nothing more to do for non-connect services
		return nil
	}

	sidecar := svc.Connect.SidecarService
	if sidecar.Proxy == nil {
		sidecar.Proxy = &nomadapi.ConsulProxy{}
	}

	net := tg.Networks[0]
	if cfg.PortLabel != "" {
		if !hasDynamicPort(tg, cfg.PortLabel) {
			net.DynamicPorts = append(net.DynamicPorts, nomadapi.Port{
				Label: cfg.PortLabel,
			})
		}

		if sidecar.Proxy.ExposeConfig == nil {
			sidecar.Proxy.ExposeConfig = &nomadapi.ConsulExposeConfig{}
		}
		expose := sidecar.Proxy.ExposeConfig

		svcPort, err := strconv.Atoi(svc.PortLabel)
		if err != nil {
			return fmt.Errorf("cannot convert port label to a number: %w", err)
		}

		expose.Path = append(expose.Path, &nomadapi.ConsulExposePath{
			Path:          cfg.Path,
			Protocol:      "http",
			LocalPathPort: svcPort,
			ListenerPort:  cfg.PortLabel,
		})
	}

	envoyPort := 9102 + i
	envoyPortLabel := fmt.Sprintf("envoy_metrics_%d", i)
	net.DynamicPorts = append(net.DynamicPorts, nomadapi.Port{
		Label: envoyPortLabel,
		To:    envoyPort,
	})

	if sidecar.Proxy.Config == nil {
		sidecar.Proxy.Config = map[string]interface{}{}
	}
	sidecar.Proxy.Config["envoy_prometheus_bind_addr"] = fmt.Sprintf("0.0.0.0:%d", envoyPort)
	svc.Meta["envoy_metrics_port"] = fmt.Sprintf("${NOMAD_HOST_PORT_%s}", envoyPortLabel)

	for j, chk := range svc.Checks {
		// Any HTTP health checks for a connect service must be exposed.
		// We generate the ports manually even though Nomad could do it because it keeps
		// the labels deterministic, so we don't end up with erroneous diffs when submitting
		// jobs.
		if chk.Type == "http" {
			portLabel := fmt.Sprintf("health_%d_%d", i, j)
			svc.Checks[j].Expose = true
			svc.Checks[j].PortLabel = portLabel
			AddPort(tg, nomadapi.Port{Label: portLabel})
		}
	}

	return nil
}

func getServiceIndex(tg *nomadapi.TaskGroup, svc *nomadapi.Service) int {
	for i, s := range tg.Services {
		if s == svc {
			return i
		}
	}
	return -1
}

func hasDynamicPort(tg *nomadapi.TaskGroup, label string) bool {
	net := tg.Networks[0]
	for _, port := range net.DynamicPorts {
		if port.Label == label {
			return true
		}
	}
	return false
}
