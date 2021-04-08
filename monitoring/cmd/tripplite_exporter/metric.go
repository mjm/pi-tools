package main

import (
	"context"
	"fmt"

	"github.com/zserge/hid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

const instrumentationName = "github.com/mjm/pi-tools/monitoring/cmd/tripplite_exporter"

func ObserveDevices(ctx context.Context, devices []hid.Device) {
	meter := metric.Must(global.Meter(instrumentationName))

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

	var (
		statusFullyDischarged        metric.Int64ValueObserver
		statusFullyCharged           metric.Int64ValueObserver
		statusBelowRemainingCapacity metric.Int64ValueObserver
		statusNeedsReplacement       metric.Int64ValueObserver
		statusDischarging            metric.Int64ValueObserver
		statusCharging               metric.Int64ValueObserver
	)

	statusObs := meter.NewBatchObserver(func(ctx context.Context, result metric.BatchObserverResult) {
		for _, d := range devices {
			statuses, err := ReadBitSet(d, ReportStatus)
			if err != nil {
				// TODO
				return
			}

			result.Observe(deviceLabels(d),
				statusBelowRemainingCapacity.Observation(boolToInt64(statuses[2])),
				statusFullyCharged.Observation(boolToInt64(statuses[3])),
				statusCharging.Observation(boolToInt64(statuses[4])),
				statusDischarging.Observation(boolToInt64(statuses[5])),
				statusFullyDischarged.Observation(boolToInt64(statuses[6])),
				statusNeedsReplacement.Observation(boolToInt64(statuses[7])))
		}
	})

	statusFullyDischarged = statusObs.NewInt64ValueObserver("tripplite.status.battery.fully_discharged")
	statusFullyCharged = statusObs.NewInt64ValueObserver("tripplite.status.battery.fully_charged")
	statusBelowRemainingCapacity = statusObs.NewInt64ValueObserver("tripplite.status.battery.below_remaining_capacity")
	statusNeedsReplacement = statusObs.NewInt64ValueObserver("tripplite.status.battery.needs_replacement")
	statusDischarging = statusObs.NewInt64ValueObserver("tripplite.status.battery.discharging")
	statusCharging = statusObs.NewInt64ValueObserver("tripplite.status.battery.charging")
}

func deviceLabels(d hid.Device) []attribute.KeyValue {
	info := d.Info()
	return []attribute.KeyValue{
		attribute.String("vendor", fmt.Sprintf("%04x", info.Vendor)),
		attribute.String("product", fmt.Sprintf("%04x", info.Product)),
	}
}

func boolToInt64(v bool) int64 {
	if v {
		return 1
	}
	return 0
}
