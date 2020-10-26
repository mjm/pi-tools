package checker

import (
	"go.opentelemetry.io/otel/api/metric"
)

type metrics struct {
	BluetoothCheckTotal       metric.Int64Counter
	BluetoothCheckErrorsTotal metric.Int64Counter

	DeviceCheckDuration    metric.Float64ValueRecorder
	DeviceCheckTotal       metric.Int64Counter
	DeviceCheckErrorsTotal metric.Int64Counter
}

func newMetrics(meter metric.Meter) metrics {
	m := metric.Must(meter)
	return metrics{
		BluetoothCheckTotal: m.NewInt64Counter("presence.bluetooth.health_check.total",
			metric.WithDescription("Counts how many times we've attempted to check the health of the local Bluetooth device")),
		BluetoothCheckErrorsTotal: m.NewInt64Counter("presence.bluetooth.health_check.errors.total",
			metric.WithDescription("Counts how many times we've failed to check the health of the local Bluetooth device")),

		DeviceCheckDuration: m.NewFloat64ValueRecorder("presence.device_check.duration.seconds",
			metric.WithDescription("Measures how long it takes to check the presence of a device")),
		DeviceCheckTotal: m.NewInt64Counter("presence.device_check.total",
			metric.WithDescription("Counts how many times we've attempted to check the presence of a device")),
		DeviceCheckErrorsTotal: m.NewInt64Counter("presence.device_check.errors.total",
			metric.WithDescription("Counts how many times we've failed to check the presence of a device")),
	}
}
