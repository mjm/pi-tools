package trips

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

var (
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
	db           *database.Client
	currentTrip  *database.Trip
	lastLeft     time.Time
	lastReturned time.Time
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
		t.setLastLeft(currentTrip.LeftAt)
	}
	if lastCompletedTrip != nil {
		t.setLastReturned(lastCompletedTrip.ReturnedAt)
		if t.lastLeft.IsZero() {
			t.setLastLeft(lastCompletedTrip.LeftAt)
		}
	}

	return t, nil
}

func (t *Tracker) OnLeave(_ *presence.Tracker) {
	t.lastLeft = time.Now()
	lastLeaveTimestamp.SetToCurrentTime()

	if _, err := t.db.BeginTrip(context.Background(), t.lastLeft); err != nil {
		log.Printf("Error saving new trip to DB: %v", err)
	}
}

func (t *Tracker) OnReturn(_ *presence.Tracker) {
	t.lastReturned = time.Now()
	lastReturnTimestamp.SetToCurrentTime()

	if !t.lastLeft.IsZero() {
		tripDuration := t.lastReturned.Sub(t.lastLeft)
		tripDurationSeconds.Observe(tripDuration.Seconds())
	}
}

func (t *Tracker) setLastLeft(ts time.Time) {
	t.lastLeft = ts
	lastLeaveTimestamp.Set(float64(ts.UnixNano()) / 1e9)
}

func (t *Tracker) setLastReturned(ts time.Time) {
	t.lastReturned = ts
	lastReturnTimestamp.Set(float64(ts.UnixNano()) / 1e9)
}
