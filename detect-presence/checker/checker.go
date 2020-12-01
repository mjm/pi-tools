package checker

import (
	"context"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/checker"

var tracer = otel.Tracer(instrumentationName)

type Checker struct {
	Tracker  *presence.Tracker
	Detector detector.Detector
	Interval time.Duration
	Devices  []presence.Device

	clock         clockwork.Clock
	metrics       metrics
	isHealthy     bool
	missingCanary bool
	lock          sync.Mutex
}

func (c *Checker) Run(ctx context.Context, tickCh chan<- struct{}) {
	if c.clock == nil {
		c.clock = clockwork.NewRealClock()
	}

	meter := otel.Meter(instrumentationName)
	c.metrics = newMetrics(meter)
	metric.Must(meter).NewInt64ValueObserver("presence.bluetooth.healthy", func(ctx context.Context, result metric.Int64ObserverResult) {
		c.lock.Lock()
		defer c.lock.Unlock()

		var val int64
		if c.isHealthy {
			val = 1
		}
		result.Observe(val)
	}, metric.WithDescription("Indicates if the local Bluetooth device is up and running."))
	metric.Must(meter).NewInt64ValueObserver("presence.canary.missing", func(ctx context.Context, result metric.Int64ObserverResult) {
		c.lock.Lock()
		defer c.lock.Unlock()

		var val int64
		if c.missingCanary {
			val = 1
		}
		result.Observe(val)
	})

	c.tick(ctx)
	if tickCh != nil {
		tickCh <- struct{}{}
	}

	ticker := c.clock.NewTicker(c.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.Chan():
			c.tick(ctx)
			if tickCh != nil {
				tickCh <- struct{}{}
			}
		}
	}
}

func (c *Checker) tick(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "Checker.tick")
	defer span.End()

	healthy := c.checkBluetoothHealth(ctx)
	span.SetAttributes(label.Bool("bluetooth.healthy", healthy))

	var presentDeviceCount, missingDeviceCount int
	var missingCanary bool
	for _, d := range c.Devices {
		if c.checkDevice(ctx, d) {
			presentDeviceCount++
		} else {
			missingDeviceCount++

			// if the missing device was a canary device, skip checking the other devices
			if d.Canary {
				missingCanary = true
				break
			}
		}
	}

	c.lock.Lock()
	c.missingCanary = missingCanary
	c.lock.Unlock()

	span.SetAttributes(
		label.Bool("device.missing_canary", missingCanary),
		label.Int("device.present_count", presentDeviceCount),
		label.Int("device.missing_count", missingDeviceCount),
		label.Int("device.count", presentDeviceCount+missingDeviceCount))
}

func (c *Checker) checkBluetoothHealth(ctx context.Context) bool {
	ctx, span := tracer.Start(ctx, "Checker.checkBluetoothHealth")
	defer span.End()

	healthy, err := c.Detector.IsHealthy(ctx)
	c.metrics.BluetoothCheckTotal.Add(ctx, 1)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.metrics.BluetoothCheckErrorsTotal.Add(ctx, 1)
	}

	span.SetAttributes(label.Bool("bluetooth.healthy", healthy))
	c.lock.Lock()
	c.isHealthy = healthy
	c.lock.Unlock()

	return healthy
}

func (c *Checker) checkDevice(ctx context.Context, d presence.Device) bool {
	ctx, span := tracer.Start(ctx, "Checker.checkDevice",
		trace.WithAttributes(
			label.String("device.name", d.Name),
			label.String("device.addr", d.Addr),
			label.Bool("device.is_canary", d.Canary)))
	defer span.End()

	startTime := time.Now()
	present, err := c.Detector.DetectDevice(ctx, d.Addr)
	duration := time.Now().Sub(startTime)

	labels := []label.KeyValue{
		label.String("name", d.Name),
		label.String("addr", d.Addr),
		label.Bool("canary", d.Canary),
	}

	c.metrics.DeviceCheckTotal.Add(ctx, 1, labels...)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.metrics.DeviceCheckErrorsTotal.Add(ctx, 1, labels...)
		return false
	}
	c.metrics.DeviceCheckDuration.Record(ctx, duration.Seconds(), labels...)

	span.SetAttributes(label.Bool("device.present", present))
	if !d.Canary {
		c.Tracker.Set(ctx, d, present)
	}
	return present
}
