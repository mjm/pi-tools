package observability

import (
	"go.opentelemetry.io/otel/sdk/trace"
)

type customSampler struct {
	defaultSampler     trace.Sampler
	healthcheckSampler trace.Sampler
	authRequestSampler trace.Sampler
}

func (s customSampler) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	if parameters.Name == "Server" {
		for _, kv := range parameters.Attributes {
			if string(kv.Key) == "http.target" {
				if kv.Value.AsString() == "/healthz" {
					return s.healthcheckSampler.ShouldSample(parameters)
				}
				if kv.Value.AsString() == "/auth" {
					for _, kv2 := range parameters.Attributes {
						if string(kv2.Key) == "http.status_code" {
							if kv.Value.AsInt64() == 200 {
								return s.authRequestSampler.ShouldSample(parameters)
							}
							break
						}
					}
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
		authRequestSampler: trace.TraceIDRatioBased(0.1),
	})
}
