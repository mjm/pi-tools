package checker

import (
	"context"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

type Checker struct {
	Tracker  *presence.Tracker
	Detector detector.Detector
	Interval time.Duration
	Devices  []presence.Device

	metrics   metrics
	isHealthy bool
	lock      sync.Mutex
}

func (c *Checker) Run() {
	meter := global.Meter("github.com/mjm/pi-tools/detect-presence/checker")
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

	c.tick()
	for range time.Tick(c.Interval) {
		c.tick()
	}
}

func (c *Checker) tick() {
	ctx := context.Background()

	// first, check the health of the bluetooth device
	healthy, err := c.Detector.IsHealthy(ctx)
	c.metrics.BluetoothCheckTotal.Add(ctx, 1)
	if err != nil {
		log.Printf("Failed to check Bluetooth device health: %v", err)
		c.metrics.BluetoothCheckErrorsTotal.Add(ctx, 1)
	}

	c.lock.Lock()
	c.isHealthy = healthy
	c.lock.Unlock()

	for _, d := range c.Devices {
		startTime := time.Now()
		present, err := c.Detector.DetectDevice(context.Background(), d.Addr)
		duration := time.Now().Sub(startTime)

		labels := []label.KeyValue{
			label.String("name", d.Name),
			label.String("addr", d.Addr),
		}

		c.metrics.DeviceCheckTotal.Add(ctx, 1, labels...)
		if err != nil {
			log.Printf("Failed to detect device %q (%s): %v", d.Name, d.Addr, err)
			c.metrics.DeviceCheckErrorsTotal.Add(ctx, 1, labels...)
			continue
		}
		c.metrics.DeviceCheckDuration.Record(ctx, duration.Seconds(), labels...)

		c.Tracker.Set(d, present)
		log.Printf("Successfully detected device %q", d.Name)
	}
}
