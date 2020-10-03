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
)

type Tracker struct {
	devices      map[Device]bool
	lastLeft     time.Time
	lastReturned time.Time
}

func NewTracker() *Tracker {
	return &Tracker{
		devices: map[Device]bool{},
	}
}

func (t *Tracker) IsPresent() bool {
	// assume present by default until we see a device that's missing
	present := true

	for _, devicePresent := range t.devices {
		// TODO make this resilient to devices being temporarily unavailable
		if !devicePresent {
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

	t.devices[d] = present

	isPresent := t.IsPresent()

	if !wasPresent && isPresent {
		// we have returned!
		log.Printf("Transitioned from away to home")
		t.lastReturned = time.Now()
		lastReturnTimestamp.SetToCurrentTime()
	} else if wasPresent && !isPresent {
		// we have abandoned our home!
		log.Printf("Transitioned from home to away")
		t.lastLeft = time.Now()
		lastLeaveTimestamp.SetToCurrentTime()
	}
}
