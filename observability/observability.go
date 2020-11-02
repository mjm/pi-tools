package observability

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagators"
	"go.opentelemetry.io/otel/semconv"

	"github.com/mjm/pi-tools/debug"
)

func Start(svcname string) (func(), error) {
	var err error
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
				ServiceName: svcname,
				Tags: []label.KeyValue{
					semconv.K8SPodNameKey.String(os.Getenv("POD_NAME")),
				},
			}))

		// TODO re-enable this once there's a jaeger agent Docker image for ARM64
		//stopTracing, err = jaeger.InstallNewPipeline(jaeger.WithAgentEndpoint("localhost:6831"))

		if err != nil {
			return nil, fmt.Errorf("installing jaeger tracing pipeline: %w", err)
		}
	}
	global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))

	// this comes after because we want the prometheus meter provider even when debugging
	metrics, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		stopTracing()
		return nil, fmt.Errorf("installing metrics pipeline: %w", err)
	}
	http.Handle("/metrics", metrics)

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
