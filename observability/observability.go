package observability

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/mjm/pi-tools/debug"
)

func Start(svcname string) (func(), error) {
	var err error
	var stopTracing func()

	if !debug.IsEnabled() {
		hostIP := os.Getenv("HOST_IP")
		exporter, err := otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:4317", hostIP)))
		if err != nil {
			return nil, fmt.Errorf("creating otlp exporter: %w", err)
		}

		r, err := resource.New(context.Background(), resource.WithAttributes(
			semconv.ServiceNamespaceKey.String(os.Getenv("NOMAD_NAMESPACE")),
			semconv.ServiceNameKey.String(svcname),
			semconv.ServiceInstanceIDKey.String(os.Getenv("NOMAD_ALLOC_ID")),

			semconv.ContainerNameKey.String(os.Getenv("NOMAD_TASK_NAME")),

			semconv.HostNameKey.String(os.Getenv("HOSTNAME")),
			semconv.HostIDKey.String(os.Getenv("NOMAD_CLIENT_ID"))))
		if err != nil {
			return nil, fmt.Errorf("creating telemetry resource: %w", err)
		}

		tp := trace.NewTracerProvider(
			trace.WithSampler(DefaultSampler()),
			trace.WithBatcher(exporter),
			trace.WithResource(r))
		otel.SetTracerProvider(tp)

		stopTracing = func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Panicf("shutting down tracing: %v", err)
			}
		}
	} else {
		stopTracing = func() {}
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// this comes after because we want the prometheus meter provider even when debugging
	ctl := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries([]float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 25, 50})),
			export.CumulativeExportKindSelector(),
			processor.WithMemory(true)))
	exporter, err := prometheus.New(prometheus.Config{}, ctl)
	if err != nil {
		stopTracing()
		return nil, fmt.Errorf("installing metrics pipeline: %w", err)
	}
	global.SetMeterProvider(exporter.MeterProvider())
	http.Handle("/metrics", exporter)

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
