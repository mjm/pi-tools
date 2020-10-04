package presence

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	deviceTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "presence",
		Name:      "device_total",
		Help:      "Indicates which devices are detected to be present at the time.",
	}, []string{"name", "addr"})
)

type Tracker struct {
	AllowedFailures int

	devices       map[Device]int
	onLeaveHooks  []OnLeaveHook
	onReturnHooks []OnReturnHook
}

func NewTracker() *Tracker {
	return &Tracker{
		AllowedFailures: 1,
		devices:         map[Device]int{},
		onLeaveHooks: []OnLeaveHook{
			loggingHook{},
		},
		onReturnHooks: []OnReturnHook{
			loggingHook{},
		},
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
		for _, hook := range t.onReturnHooks {
			hook.OnReturn(t)
		}
	} else if wasPresent && !isPresent {
		// we have abandoned our home!
		for _, hook := range t.onLeaveHooks {
			hook.OnLeave(t)
		}
	}
}
