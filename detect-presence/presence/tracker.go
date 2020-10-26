package presence

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/label"
)

type Tracker struct {
	AllowedFailures int

	devices       map[Device]int
	onLeaveHooks  []OnLeaveHook
	onReturnHooks []OnReturnHook
	lock          sync.Mutex
}

func NewTracker() *Tracker {
	t := &Tracker{
		AllowedFailures: 2,
		devices:         map[Device]int{},
		onLeaveHooks: []OnLeaveHook{
			loggingHook{},
		},
		onReturnHooks: []OnReturnHook{
			loggingHook{},
		},
	}

	m := metric.Must(global.Meter("github.com/mjm/pi-tools/detect-presence/presence"))
	m.NewInt64ValueObserver("presence.device.total", func(ctx context.Context, result metric.Int64ObserverResult) {
		t.lock.Lock()
		defer t.lock.Unlock()

		for d, n := range t.devices {
			var val int64
			if n == 0 {
				val = 1
			}

			result.Observe(val, label.String("name", d.Name), label.String("addr", d.Addr))
		}
	}, metric.WithDescription("Indicates which devices are detected to be present at the time."))

	return t
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
	t.lock.Lock()
	wasPresent := t.IsPresent()

	if present {
		t.devices[d] = 0
	} else {
		t.devices[d] += 1
	}

	isPresent := t.IsPresent()
	t.lock.Unlock()

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
