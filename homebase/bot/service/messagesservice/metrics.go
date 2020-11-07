package messagesservice

import (
	"go.opentelemetry.io/otel/api/metric"
)

type metrics struct {
	TelegramUpdateTotal       metric.Int64Counter
	TelegramUpdateErrorsTotal metric.Int64Counter
}

func newMetrics(meter metric.Meter) metrics {
	m := metric.Must(meter)
	return metrics{
		TelegramUpdateTotal: m.NewInt64Counter("homebase.bot.telegram.update.total",
			metric.WithDescription("Counts how many individual updates have been received from Telegram")),
		TelegramUpdateErrorsTotal: m.NewInt64Counter("homebase.bot.telegram.update.errors.total",
			metric.WithDescription("Counts how many errors have been receive from trying to poll for updates from Telegram")),
	}
}
