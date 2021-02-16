package observability

import (
	"go.opentelemetry.io/otel/sdk/trace"
)

type customSampler struct {
	defaultSampler     trace.Sampler
	healthcheckSampler trace.Sampler
}

func (s customSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	if parameters.Name == "Server" {
		for _, kv := range parameters.Attributes {
			if string(kv.Key) == "http.target" {
				if kv.Value.AsString() == "/healthz" {
					return s.healthcheckSampler.ShouldSample(parameters)
				}
				break
			}
		}
	}

	return s.defaultSampler.ShouldSample(parameters)
}

func (s customSampler) Description() string {
	return "PiToolsSampler"
}

func DefaultSampler() trace.Sampler {
	return trace.ParentBased(customSampler{
		defaultSampler:     trace.AlwaysSample(),
		healthcheckSampler: trace.TraceIDRatioBased(0.1),
	})
}
