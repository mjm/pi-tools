package trips

import (
	"context"
	"log"
	"sync"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

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

	m := metric.Must(global.Meter("github.com/mjm/pi-tools/detect-presence/trips"))
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

func (t *Tracker) OnLeave(_ *presence.Tracker) {
	if t.currentTrip != nil {
		log.Printf("Skipping starting new trip because there's already a current trip")
		return
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastLeft = time.Now()

	newTrip, err := t.db.BeginTrip(context.Background(), t.lastLeft)
	if err != nil {
		log.Printf("Error saving new trip to DB: %v", err)
	} else {
		t.currentTrip = newTrip
	}
}

func (t *Tracker) OnReturn(_ *presence.Tracker) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastReturned = time.Now()

	if !t.lastLeft.IsZero() {
		tripDuration := t.lastReturned.Sub(t.lastLeft)
		t.tripDurationSeconds.Record(context.Background(), tripDuration.Seconds())
	}

	if t.currentTrip != nil {
		if err := t.db.EndTrip(context.Background(), t.currentTrip.ID, t.lastReturned); err != nil {
			log.Printf("Error completing trip in DB: %v", err)
		}
		t.currentTrip = nil
	}
}
