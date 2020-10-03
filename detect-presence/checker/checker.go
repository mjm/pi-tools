package checker

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

var (
	bluetoothHealthy = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "presence",
		Name:      "bluetooth_healthy",
		Help:      "Indicates if the local Bluetooth device is up and running.",
	})

	deviceCheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "presence",
		Name:      "device_check_duration_seconds",
		Help:      "Measures how long it takes to check the presence of a device",
		Buckets:   []float64{0.25, 0.5, 1.0, 2.0, 3.0, 5.0, 10.0, 30.0},
	}, []string{"name", "addr"})
)

type Checker struct {
	Tracker    *presence.Tracker
	Interval   time.Duration
	DeviceName string
	Devices    []presence.Device
}

func (c *Checker) Run() {
	c.tick()
	for range time.Tick(c.Interval) {
		c.tick()
	}
}

func (c *Checker) tick() {
	// first, check the health of the bluetooth device
	healthy, err := detector.IsHealthy(context.Background(), c.DeviceName)
	if err != nil {
		log.Printf("Failed to check Bluetooth device %q health: %v", c.DeviceName, err)
	}
	var healthyVal float64
	if healthy {
		healthyVal = 1.0
	}
	bluetoothHealthy.Set(healthyVal)

	for _, d := range c.Devices {
		startTime := time.Now()
		present, err := detector.DetectDevice(context.Background(), d.Addr)
		if err != nil {
			log.Printf("Failed to detect device %q (%s): %v", d.Name, d.Addr, err)
			continue
		}
		duration := time.Now().Sub(startTime)
		deviceCheckDuration.WithLabelValues(d.Name, d.Addr).Observe(duration.Seconds())

		c.Tracker.Set(d, present)
		log.Printf("Successfully detected device %q", d.Name)
	}
}
