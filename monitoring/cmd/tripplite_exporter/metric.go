package main

import (
	"context"
	"fmt"

	"github.com/zserge/hid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
)

const instrumentationName = "github.com/mjm/pi-tools/monitoring/cmd/tripplite_exporter"

func ObserveDevices(ctx context.Context, devices []hid.Device) {
	meter := metric.Must(otel.Meter(instrumentationName))

	meter.NewInt64ValueObserver("tripplite.health", func(ctx context.Context, result metric.Int64ObserverResult) {
		for _, d := range devices {
			health, err := ReadInt8(d, ReportHealth)
			if err != nil {
				// TODO
				return
			}
			result.Observe(int64(health), deviceLabels(d)...)
		}
	}, metric.WithDescription("The percentage of charge remaining in the battery"))

	meter.NewInt64ValueObserver("tripplite.config.voltage", func(ctx context.Context, result metric.Int64ObserverResult) {
		for _, d := range devices {
			voltage, err := ReadInt8(d, ReportConfigVoltage)
			if err != nil {
				// TODO
				return
			}
			result.Observe(int64(voltage), deviceLabels(d)...)
		}
	})

	meter.NewInt64ValueObserver("tripplite.config.frequency", func(ctx context.Context, result metric.Int64ObserverResult) {
		for _, d := range devices {
			freq, err := ReadInt8(d, ReportConfigFrequency)
			if err != nil {
				// TODO
				return
			}
			result.Observe(int64(freq), deviceLabels(d)...)
		}
	})

	meter.NewInt64ValueObserver("tripplite.config.power", func(ctx context.Context, result metric.Int64ObserverResult) {
		for _, d := range devices {
			freq, err := ReadInt16(d, ReportConfigPower)
			if err != nil {
				// TODO
				return
			}
			result.Observe(int64(freq), deviceLabels(d)...)
		}
	})

	meter.NewFloat64ValueObserver("tripplite.input.voltage", func(ctx context.Context, result metric.Float64ObserverResult) {
		for _, d := range devices {
			voltage, err := ReadFloat(d, ReportInputVoltage)
			if err != nil {
				// TODO
				return
			}
			result.Observe(voltage, deviceLabels(d)...)
		}
	})

	meter.NewFloat64ValueObserver("tripplite.input.frequency", func(ctx context.Context, result metric.Float64ObserverResult) {
		for _, d := range devices {
			voltage, err := ReadFloat(d, ReportInputFrequency)
			if err != nil {
				// TODO
				return
			}
			result.Observe(voltage, deviceLabels(d)...)
		}
	})

	meter.NewFloat64ValueObserver("tripplite.output.voltage", func(ctx context.Context, result metric.Float64ObserverResult) {
		for _, d := range devices {
			voltage, err := ReadFloat(d, ReportOutputVoltage)
			if err != nil {
				// TODO
				return
			}
			result.Observe(voltage, deviceLabels(d)...)
		}
	})

	meter.NewInt64ValueObserver("tripplite.output.power", func(ctx context.Context, result metric.Int64ObserverResult) {
		for _, d := range devices {
			freq, err := ReadInt16(d, ReportOutputPower)
			if err != nil {
				// TODO
				return
			}
			result.Observe(int64(freq), deviceLabels(d)...)
		}
	})

	meter.NewInt64ValueObserver("tripplite.time_to_empty.seconds", func(ctx context.Context, result metric.Int64ObserverResult) {
		for _, d := range devices {
			freq, err := ReadInt16(d, ReportTimeToEmpty)
			if err != nil {
				// TODO
				return
			}
			result.Observe(int64(freq), deviceLabels(d)...)
		}
	})
}

func deviceLabels(d hid.Device) []label.KeyValue {
	info := d.Info()
	return []label.KeyValue{
		label.String("vendor", fmt.Sprintf("%04x", info.Vendor)),
		label.String("product", fmt.Sprintf("%04x", info.Product)),
	}
}
