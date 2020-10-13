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

	bluetoothCheckTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "presence",
		Name:      "bluetooth_health_check_total",
		Help:      "Counts how many times we've attempted to check the health of the local Bluetooth device",
	})

	bluetoothCheckErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "presence",
		Name:      "bluetooth_health_check_errors_total",
		Help:      "Counts how many times we've failed to check the health of the local Bluetooth device",
	})

	deviceCheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "presence",
		Name:      "device_check_duration_seconds",
		Help:      "Measures how long it takes to check the presence of a device",
		Buckets:   []float64{0.25, 0.5, 1.0, 2.0, 3.0, 5.0, 10.0, 30.0},
	}, []string{"name", "addr"})

	deviceCheckTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "presence",
		Name:      "device_check_total",
		Help:      "Counts how many times we've attempted to check the presence of a device",
	}, []string{"name", "addr"})

	deviceCheckErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "presence",
		Name:      "device_check_errors_total",
		Help:      "Counts how many times we've failed to check the presence of a device",
	}, []string{"name", "addr"})
)

type Checker struct {
	Tracker    *presence.Tracker
	Detector   detector.Detector
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
	healthy, err := c.Detector.IsHealthy(context.Background(), c.DeviceName)
	bluetoothCheckTotal.Inc()
	if err != nil {
		log.Printf("Failed to check Bluetooth device %q health: %v", c.DeviceName, err)
		bluetoothCheckErrorsTotal.Inc()
	}
	var healthyVal float64
	if healthy {
		healthyVal = 1.0
	}
	bluetoothHealthy.Set(healthyVal)

	for _, d := range c.Devices {
		startTime := time.Now()
		present, err := c.Detector.DetectDevice(context.Background(), d.Addr)
		duration := time.Now().Sub(startTime)
		deviceCheckTotal.WithLabelValues(d.Name, d.Addr).Inc()
		if err != nil {
			log.Printf("Failed to detect device %q (%s): %v", d.Name, d.Addr, err)
			deviceCheckErrorsTotal.WithLabelValues(d.Name, d.Addr).Inc()
			continue
		}
		deviceCheckDuration.WithLabelValues(d.Name, d.Addr).Observe(duration.Seconds())

		c.Tracker.Set(d, present)
		log.Printf("Successfully detected device %q", d.Name)
	}
}
