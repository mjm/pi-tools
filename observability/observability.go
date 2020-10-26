package observability

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"

	"github.com/mjm/pi-tools/debug"
)

func Start() (func(), error) {
	metrics, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil, fmt.Errorf("installing metrics pipeline: %w", err)
	}
	http.Handle("/metrics", metrics)

	var stopTracing func()
	if debug.IsEnabled() {
		tracePipe, err := stdout.InstallNewPipeline([]stdout.Option{stdout.WithPrettyPrint()}, nil)
		if err != nil {
			log.Panicf("Error installing stdout tracing pipeline: %v", err)
		}
		stopTracing = tracePipe.Stop
	} else {
		stopTracing, err = jaeger.InstallNewPipeline(
			jaeger.WithCollectorEndpoint("http://jaeger-collector.monitoring:14268/api/traces"),
			jaeger.WithProcess(jaeger.Process{
				ServiceName: "go-links",
				Tags: []label.KeyValue{
					semconv.K8SPodNameKey.String(os.Getenv("POD_NAME")),
				},
			}),
			jaeger.WithSDK(&sdktrace.Config{
				DefaultSampler: noMetricsSampler{},
			}))

		// TODO re-enable this once there's a jaeger agent Docker image for ARM64
		//stopTracing, err = jaeger.InstallNewPipeline(jaeger.WithAgentEndpoint("localhost:6831"))

		if err != nil {
			return nil, fmt.Errorf("installing jaeger tracing pipeline: %w", err)
		}
	}

	return stopTracing, nil
}
