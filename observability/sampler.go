package observability

import (
	"go.opentelemetry.io/otel/api/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type noMetricsSampler struct{}

func (noMetricsSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if p.Kind != trace.SpanKindServer {
		return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}
	}

	for _, attr := range p.Attributes {
		if attr.Key == "http.target" && attr.Value.AsString() == "/metrics" {
			return sdktrace.SamplingResult{Decision: sdktrace.Drop}
		}
	}

	return sdktrace.SamplingResult{Decision: sdktrace.RecordAndSample}
}

func (noMetricsSampler) Description() string {
	return "NoMetricsSampler"
}
