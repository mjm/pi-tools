package trips

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/presence"
	messagespb "github.com/mjm/pi-tools/homebase/bot/proto/messages"
	"github.com/mjm/pi-tools/pkg/spanerr"
	"github.com/mjm/pi-tools/storage"
)

const instrumentationName = "github.com/mjm/pi-tools/detect-presence/trips"

var tracer = global.Tracer(instrumentationName)

type Tracker struct {
	db                  *database.Queries
	messages            messagespb.MessagesServiceClient
	clock               clockwork.Clock
	currentTrip         *database.Trip
	lastLeft            time.Time
	lastReturned        time.Time
	tripDurationSeconds metric.Float64ValueRecorder
	lock                sync.Mutex
}

func NewTracker(db storage.DB, messages messagespb.MessagesServiceClient) (*Tracker, error) {
	ctx := context.Background()
	q := database.New(db)

	var missingLastCompleted, missingCurrent bool
	lastCompletedTrip, err := q.GetLastCompletedTrip(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			missingLastCompleted = true
		} else {
			return nil, err
		}
	}

	currentTrip, err := q.GetCurrentTrip(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			missingCurrent = true
		} else {
			return nil, err
		}
	}

	t := &Tracker{
		db:       q,
		clock:    clockwork.NewRealClock(),
		messages: messages,
	}

	// re-populate state and metrics to pick up where a previous process left off
	if !missingCurrent {
		t.currentTrip = &currentTrip
		t.lastLeft = currentTrip.LeftAt
	}
	if !missingLastCompleted {
		t.lastReturned = lastCompletedTrip.ReturnedAt.Time
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

	t.lastLeft = t.clock.Now()
	span.SetAttributes(label.String("trip.left_at", t.lastLeft.UTC().Format(time.RFC3339)))

	id := uuid.New()
	span.SetAttributes(label.String("trip.id", id.String()))

	newTrip, err := t.db.BeginTrip(ctx, database.BeginTripParams{
		ID:     id,
		LeftAt: t.lastLeft,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		t.currentTrip = &newTrip
	}
}

func (t *Tracker) OnReturn(ctx context.Context, _ *presence.Tracker) {
	ctx, span := tracer.Start(ctx, "trips.Tracker.OnLeave",
		trace.WithAttributes(label.Bool("trip.in_progress", t.currentTrip != nil)))
	defer span.End()

	t.lock.Lock()
	defer t.lock.Unlock()

	t.lastReturned = t.clock.Now()
	span.SetAttributes(
		label.String("trip.returned_at", t.lastReturned.UTC().Format(time.RFC3339)),
		label.String("trip.left_at", t.lastLeft.UTC().Format(time.RFC3339)))

	if !t.lastLeft.IsZero() {
		tripDuration := t.lastReturned.Sub(t.lastLeft)
		span.SetAttributes(label.Float64("trip.duration_secs", tripDuration.Seconds()))
		t.tripDurationSeconds.Record(ctx, tripDuration.Seconds())
	}

	if t.currentTrip != nil {
		span.SetAttributes(label.String("trip.id", t.currentTrip.ID.String()))
		if err := t.db.EndTrip(ctx, database.EndTripParams{
			ID: t.currentTrip.ID,
			ReturnedAt: sql.NullTime{
				Time:  t.lastReturned,
				Valid: true,
			},
		}); err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		t.sendChatMessage(ctx)

		t.currentTrip = nil
	}
}

func (t *Tracker) sendChatMessage(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "trips.Tracker.sendChatMessage")
	defer span.End()

	req := &messagespb.SendTripCompletedMessageRequest{
		TripId:     t.currentTrip.ID.String(),
		LeftAt:     t.lastLeft.Format(time.RFC3339),
		ReturnedAt: t.lastReturned.Format(time.RFC3339),
	}
	if _, err := t.messages.SendTripCompletedMessage(ctx, req); err != nil {
		spanerr.RecordError(ctx, err)
	}
}
