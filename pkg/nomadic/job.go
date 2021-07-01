package nomadic

import (
	"fmt"
	"strconv"

	nomadapi "github.com/hashicorp/nomad/api"
)

func NewJob(name string, priority int) *nomadapi.Job {
	return &nomadapi.Job{
		ID:          &name,
		Datacenters: DefaultDatacenters,
		Priority:    &priority,
	}
}

func NewSystemJob(name string, priority int) *nomadapi.Job {
	j := NewJob(name, priority)
	j.Type = String(nomadapi.JobTypeSystem)
	return j
}

func NewBatchJob(name string, priority int) *nomadapi.Job {
	j := NewJob(name, priority)
	j.Type = String(nomadapi.JobTypeBatch)
	return j
}

func AddTaskGroup(job *nomadapi.Job, name string, count int) *nomadapi.TaskGroup {
	tg := &nomadapi.TaskGroup{
		Name:  &name,
		Count: &count,
		Networks: []*nomadapi.NetworkResource{
			{
				DNS: DefaultDNS,
			},
		},
	}
	job.AddTaskGroup(tg)
	return tg
}

type ServiceOption func(tg *nomadapi.TaskGroup, svc *nomadapi.Service)

func WithMetricsScraping(path string) ServiceOption {
	return func(tg *nomadapi.TaskGroup, svc *nomadapi.Service) {
		svc.Meta["metrics_path"] = path
		if svc.Connect == nil {
			return
		}

		AddPort(tg, nomadapi.Port{Label: "expose"})
		svc.Meta["metrics_port"] = "${NOMAD_HOST_PORT_expose}"

		svcPort, err := strconv.Atoi(svc.PortLabel)
		if err != nil {
			panic(fmt.Errorf("cannot convert port label to a number: %w", err))
		}

		sidecar := svc.Connect.SidecarService
		if sidecar.Proxy.ExposeConfig == nil {
			sidecar.Proxy.ExposeConfig = &nomadapi.ConsulExposeConfig{}
		}
		expose := sidecar.Proxy.ExposeConfig

		expose.Path = append(expose.Path, &nomadapi.ConsulExposePath{
			Path:          path,
			Protocol:      "http",
			LocalPathPort: svcPort,
			ListenerPort:  "expose",
		})
	}
}

func WithMetricsPort(label string) ServiceOption {
	return func(tg *nomadapi.TaskGroup, svc *nomadapi.Service) {
		if svc.Connect != nil {
			panic(fmt.Errorf("overriding metrics port for a connect service is not allowed"))
		}

		svc.Meta["metrics_port"] = fmt.Sprintf("${NOMAD_HOST_PORT_%s}", label)
	}
}

func WithUpstreams(upstreams ...*nomadapi.ConsulUpstream) ServiceOption {
	return func(tg *nomadapi.TaskGroup, svc *nomadapi.Service) {
		proxy := svc.Connect.SidecarService.Proxy
		proxy.Upstreams = append(proxy.Upstreams, upstreams...)
	}
}

func AddConnectService(tg *nomadapi.TaskGroup, svc *nomadapi.Service, opts ...ServiceOption) *nomadapi.Service {
	i := len(tg.Services)

	if svc.Meta == nil {
		svc.Meta = map[string]string{}
	}
	if svc.Connect == nil {
		svc.Connect = &nomadapi.ConsulConnect{}
	}
	if svc.Connect.SidecarService == nil {
		svc.Connect.SidecarService = &nomadapi.ConsulSidecarService{}
	}
	sidecar := svc.Connect.SidecarService
	if sidecar.Proxy == nil {
		sidecar.Proxy = &nomadapi.ConsulProxy{}
	}

	tg.Networks[0].Mode = "bridge"
	tg.Services = append(tg.Services, svc)

	for _, opt := range opts {
		opt(tg, svc)
	}

	envoyPort := 9102 + i
	envoyPortLabel := fmt.Sprintf("envoy_metrics_%d", i)
	AddPort(tg, nomadapi.Port{
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

	return svc
}

func AddService(tg *nomadapi.TaskGroup, svc *nomadapi.Service, opts ...ServiceOption) *nomadapi.Service {
	if svc.Meta == nil {
		svc.Meta = map[string]string{}
	}

	tg.Services = append(tg.Services, svc)

	for _, opt := range opts {
		opt(tg, svc)
	}

	return svc
}

func AddPort(tg *nomadapi.TaskGroup, port nomadapi.Port) {
	if port.Value == 0 {
		if hasDynamicPort(tg, port.Label) {
			return
		}

		tg.Networks[0].DynamicPorts = append(tg.Networks[0].DynamicPorts, port)
	} else {
		// TODO check if it's already there
		tg.Networks[0].ReservedPorts = append(tg.Networks[0].ReservedPorts, port)
	}
}

type TaskOption func(tg *nomadapi.TaskGroup, task *nomadapi.Task)

func WithCPU(cpu int) TaskOption {
	return func(tg *nomadapi.TaskGroup, task *nomadapi.Task) {
		if task.Resources == nil {
			task.Resources = &nomadapi.Resources{}
		}
		task.Resources.CPU = &cpu
	}
}

func WithMemoryMB(memory int) TaskOption {
	return func(tg *nomadapi.TaskGroup, task *nomadapi.Task) {
		if task.Resources == nil {
			task.Resources = &nomadapi.Resources{}
		}
		task.Resources.MemoryMB = &memory
	}
}

func WithLoggingTag(tag string) TaskOption {
	return func(tg *nomadapi.TaskGroup, task *nomadapi.Task) {
		task.SetConfig("logging", Logging(tag))
	}
}

func WithVaultPolicies(policies ...string) TaskOption {
	return func(tg *nomadapi.TaskGroup, task *nomadapi.Task) {
		if task.Vault == nil {
			task.Vault = &nomadapi.Vault{}
		}
		task.Vault.Policies = append(task.Vault.Policies, policies...)
	}
}

func WithVaultChangeNoop() TaskOption {
	return func(tg *nomadapi.TaskGroup, task *nomadapi.Task) {
		if task.Vault == nil {
			task.Vault = &nomadapi.Vault{}
		}
		task.Vault.ChangeMode = String("noop")
	}
}

func WithTracingEnv() TaskOption {
	return func(tg *nomadapi.TaskGroup, task *nomadapi.Task) {
		if task.Env == nil {
			task.Env = map[string]string{}
		}
		task.Env["HOSTNAME"] = "${attr.unique.hostname}"
		task.Env["HOST_IP"] = "${attr.unique.network.ip-address}"
		task.Env["NOMAD_CLIENT_ID"] = "${node.unique.id}"
	}
}

func AddTask(tg *nomadapi.TaskGroup, task *nomadapi.Task, opts ...TaskOption) *nomadapi.Task {
	if task.Driver == "" {
		task.Driver = "docker"
	}

	for _, opt := range opts {
		opt(tg, task)
	}

	tg.AddTask(task)
	return task
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
