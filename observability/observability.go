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
	"go.opentelemetry.io/otel/semconv"

	"github.com/mjm/pi-tools/debug"
)

func Start(svcname string) (func(), error) {
	var err error
	var stopTracing func()

	var endpoint jaeger.EndpointOption
	if debug.IsEnabled() {
		endpoint = jaeger.WithAgentEndpoint("localhost:6831")
	} else {
		endpoint = jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces")
	}
	stopTracing, err = jaeger.InstallNewPipeline(
		endpoint,
		jaeger.WithProcess(jaeger.Process{
			ServiceName: svcname,
			Tags: []label.KeyValue{
				semconv.K8SPodNameKey.String(os.Getenv("POD_NAME")),
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
