package telegram

import (
	"go.opentelemetry.io/otel/api/metric"
)

type metrics struct {
	RequestTotal           metric.Int64Counter
	RequestErrorsTotal     metric.Int64Counter
	RequestDurationSeconds metric.Float64ValueRecorder
}

func newMetrics(meter metric.Meter) metrics {
	m := metric.Must(meter)
	return metrics{
		RequestTotal: m.NewInt64Counter("telegram.request.total",
			metric.WithDescription("Counts the number of HTTP requests made to Telegram")),
		RequestErrorsTotal: m.NewInt64Counter("telegram.request.errors.total",
			metric.WithDescription("Counts the number of HTTP requests made to Telegram that resulted in an error")),
		RequestDurationSeconds: m.NewFloat64ValueRecorder("telegram.request.duration.seconds",
			metric.WithDescription("Measures the duration of requests made to Telegram")),
	}
}
