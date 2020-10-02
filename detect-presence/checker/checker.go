package checker

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/mjm/pi-tools/detect-presence/detector"
)

var (
	deviceTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "presence",
		Name:      "device_total",
		Help:      "Indicates which devices are detected to be present at the time.",
	}, []string{"name", "addr"})

	deviceCheckDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "presence",
		Name:      "device_check_duration_seconds",
		Help:      "Measures how long it takes to check the presence of a device",
		Buckets:   []float64{0.25, 0.5, 1.0, 2.0, 3.0, 5.0, 10.0, 30.0},
	}, []string{"name", "addr"})
)

type Checker struct {
	Interval time.Duration
	Devices  []Device
}

func (c *Checker) Run() {
	c.tick()
	for range time.Tick(c.Interval) {
		c.tick()
	}
}

func (c *Checker) tick() {
	for _, d := range c.Devices {
		startTime := time.Now()
		present, err := detector.DetectDevice(context.Background(), d.Addr)
		if err != nil {
			log.Printf("Failed to detect device %q (%s): %v", d.Name, d.Addr, err)
			continue
		}
		duration := time.Now().Sub(startTime)

		var val float64
		if present {
			val = 1.0
		}
		deviceTotal.WithLabelValues(d.Name, d.Addr).Set(val)
		deviceCheckDuration.WithLabelValues(d.Name, d.Addr).Observe(duration.Seconds())

		log.Printf("Successfully detected device %q", d.Name)
	}
}
