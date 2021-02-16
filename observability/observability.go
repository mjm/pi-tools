package observability

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"

	"github.com/mjm/pi-tools/debug"
)

func Start(svcname string) (func(), error) {
	var err error
	var stopTracing func()

	var endpoint jaeger.EndpointOption
	if debug.IsEnabled() {
		endpoint = jaeger.WithAgentEndpoint("127.0.0.1:6831")
	} else {
		endpoint = jaeger.WithCollectorEndpoint("http://127.0.0.1:14268/api/traces")
	}
	stopTracing, err = jaeger.InstallNewPipeline(
		endpoint,
		jaeger.WithSDK(&trace.Config{
			DefaultSampler: DefaultSampler(),
		}),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: svcname,
			Tags: []label.KeyValue{
				semconv.ServiceNamespaceKey.String(os.Getenv("NOMAD_NAMESPACE")),
				semconv.ServiceNameKey.String(fmt.Sprintf("%s/%s",
					os.Getenv("NOMAD_JOB_NAME"),
					os.Getenv("NOMAD_GROUP_NAME"))),
				semconv.ServiceInstanceIDKey.String(os.Getenv("NOMAD_ALLOC_ID")),

				semconv.ContainerNameKey.String(os.Getenv("NOMAD_TASK_NAME")),

				semconv.HostNameKey.String(os.Getenv("HOSTNAME")),
				semconv.HostIDKey.String(os.Getenv("NOMAD_CLIENT_ID")),
			},
		}))

	if err != nil {
		return nil, fmt.Errorf("installing jaeger tracing pipeline: %w", err)
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// this comes after because we want the prometheus meter provider even when debugging
	metrics, err := prometheus.InstallNewPipeline(prometheus.Config{
		DefaultHistogramBoundaries: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 25, 50},
	})
	if err != nil {
		stopTracing()
		return nil, fmt.Errorf("installing metrics pipeline: %w", err)
	}
	http.Handle("/metrics", metrics)

	if err := runtime.Start(); err != nil {
		stopTracing()
		return nil, fmt.Errorf("starting observing runtime metrics: %w", err)
	}

	return func() {
		log.Printf("Shutting down observability...")
		stopTracing()
	}, nil
}

func MustStart(svcname string) func() {
	stop, err := Start(svcname)
	if err != nil {
		log.Panicf("Error setting up observability: %v", err)
	}
	return stop
}
