package presence

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/presence"

var tracer = otel.Tracer(instrumentationName)

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

	m := metric.Must(otel.Meter(instrumentationName))
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

func (t *Tracker) Set(ctx context.Context, d Device, present bool) {
	ctx, span := tracer.Start(ctx, "presence.Tracker.Set",
		trace.WithAttributes(
			label.String("device.name", d.Name),
			label.String("device.addr", d.Addr),
			label.Bool("device.present", present)))
	defer span.End()

	t.lock.Lock()
	wasPresent := t.IsPresent()
	span.SetAttributes(label.Bool("user.present.previous", wasPresent))

	if present {
		t.devices[d] = 0
	} else {
		t.devices[d] += 1
	}
	span.SetAttributes(label.Int("device.failure_count", t.devices[d]))

	isPresent := t.IsPresent()
	t.lock.Unlock()

	span.SetAttributes(label.Bool("user.present", isPresent))

	if !wasPresent && isPresent {
		// we have returned!
		for _, hook := range t.onReturnHooks {
			hook.OnReturn(ctx, t)
		}
	} else if wasPresent && !isPresent {
		// we have abandoned our home!
		for _, hook := range t.onLeaveHooks {
			hook.OnLeave(ctx, t)
		}
	}
}
