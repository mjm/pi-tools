package trips

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/trips"

var tracer = global.Tracer(instrumentationName)

type Tracker struct {
	db                  *database.Client
	currentTrip         *database.Trip
	lastLeft            time.Time
	lastReturned        time.Time
	tripDurationSeconds metric.Float64ValueRecorder
	lock                sync.Mutex
}

func NewTracker(db *database.Client) (*Tracker, error) {
	ctx := context.Background()
	lastCompletedTrip, err := db.GetLastCompletedTrip(ctx)
	if err != nil {
		return nil, err
	}

	currentTrip, err := db.GetCurrentTrip(ctx)
	if err != nil {
		return nil, err
	}

	t := &Tracker{
		db:          db,
		currentTrip: currentTrip,
	}

	// re-populate state and metrics to pick up where a previous process left off
	if currentTrip != nil {
		t.lastLeft = currentTrip.LeftAt
	}
	if lastCompletedTrip != nil {
		t.lastReturned = lastCompletedTrip.ReturnedAt
		if t.lastLeft.IsZero() {
			t.lastLeft = lastCompletedTrip.LeftAt
		}
	}

	m := metric.Must(global.Meter(instrumentationName))
	t.tripDurationSeconds = m.NewFloat64ValueRecorder("presence.trip.duration.seconds",
		metric.WithDescription("Measures how long trips away from home last"))

	m.NewFloat64ValueObserver("presence.last_leave.timestamp", func(ctx context.Context, result metric.Float64ObserverResult) {
		t.lock.Lock()
		defer t.lock.Unlock()

		result.Observe(float64(t.lastLeft.UnixNano()) / 1e9)
	}, metric.WithDescription("Tracks the timestamp when we last left the home."))

	m.NewFloat64ValueObserver("presence.last_return.timestamp", func(ctx context.Context, result metric.Float64ObserverResult) {
		t.lock.Lock()
		defer t.lock.Unlock()

		result.Observe(float64(t.lastReturned.UnixNano()) / 1e9)
	}, metric.WithDescription("Tracks the timestamp when we last returned to the home."))

	return t, nil
}

func (t *Tracker) OnLeave(ctx context.Context, _ *presence.Tracker) {
	ctx, span := tracer.Start(ctx, "trips.Tracker.OnLeave",
		trace.WithAttributes(label.Bool("trip.in_progress", t.currentTrip != nil)))
	defer span.End()

	if t.currentTrip != nil {
		return
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastLeft = time.Now()
	span.SetAttributes(label.String("trip.left_at", t.lastLeft.UTC().Format(time.RFC3339)))

	newTrip, err := t.db.BeginTrip(ctx, t.lastLeft)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		t.currentTrip = newTrip
		span.SetAttributes(label.String("trip.id", t.currentTrip.ID))
	}
}

func (t *Tracker) OnReturn(ctx context.Context, _ *presence.Tracker) {
	ctx, span := tracer.Start(ctx, "trips.Tracker.OnLeave",
		trace.WithAttributes(label.Bool("trip.in_progress", t.currentTrip != nil)))
	defer span.End()

	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastReturned = time.Now()
	span.SetAttributes(
		label.String("trip.returned_at", t.lastReturned.UTC().Format(time.RFC3339)),
		label.String("trip.left_at", t.lastLeft.UTC().Format(time.RFC3339)))

	if !t.lastLeft.IsZero() {
		tripDuration := t.lastReturned.Sub(t.lastLeft)
		span.SetAttributes(label.Float64("trip.duration_secs", tripDuration.Seconds()))
		t.tripDurationSeconds.Record(ctx, tripDuration.Seconds())
	}

	if t.currentTrip != nil {
		span.SetAttributes(label.String("trip.id", t.currentTrip.ID))
		if err := t.db.EndTrip(ctx, t.currentTrip.ID, t.lastReturned); err != nil {
			span.SetStatus(codes.Error, err.Error())
		}
		t.currentTrip = nil
	}
}
