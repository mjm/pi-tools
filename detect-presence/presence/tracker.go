package presence

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	deviceTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "presence",
		Name:      "device_total",
		Help:      "Indicates which devices are detected to be present at the time.",
	}, []string{"name", "addr"})

	lastLeaveTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "presence",
		Name:      "last_leave_timestamp",
		Help:      "Tracks the timestamp when we last left the home.",
	})

	lastReturnTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "presence",
		Name:      "last_return_timestamp",
		Help:      "Tracks the timestamp when we last returned to the home.",
	})

	tripDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "presence",
		Name:      "trip_duration_seconds",
		Help:      "Measures how long trips away from home last",
		Buckets:   []float64{30, 60, 180, 300, 600, 1800, 3600, 14400, 86400},
	})
)

type Tracker struct {
	AllowedFailures int

	devices      map[Device]int
	lastLeft     time.Time
	lastReturned time.Time
}

func NewTracker() *Tracker {
	return &Tracker{
		AllowedFailures: 1,
		devices:         map[Device]int{},
	}
}

func (t *Tracker) IsPresent() bool {
	// assume present by default until we see a device that's missing
	present := true

	for _, missingCount := range t.devices {
		// a device must be missing at least two pings in a row for us to be considered not present
		if missingCount > t.AllowedFailures {
			present = false
		}
	}

	return present
}

func (t *Tracker) Set(d Device, present bool) {
	wasPresent := t.IsPresent()

	var val float64
	if present {
		val = 1.0
	}
	deviceTotal.WithLabelValues(d.Name, d.Addr).Set(val)

	if present {
		t.devices[d] = 0
	} else {
		t.devices[d] += 1
	}

	isPresent := t.IsPresent()

	if !wasPresent && isPresent {
		// we have returned!
		log.Printf("Transitioned from away to home")
		t.lastReturned = time.Now()
		lastReturnTimestamp.SetToCurrentTime()

		if !t.lastLeft.IsZero() {
			tripDuration := t.lastReturned.Sub(t.lastLeft)
			tripDurationSeconds.Observe(tripDuration.Seconds())
		}
	} else if wasPresent && !isPresent {
		// we have abandoned our home!
		log.Printf("Transitioned from home to away")
		t.lastLeft = time.Now()
		lastLeaveTimestamp.SetToCurrentTime()
	}
}
