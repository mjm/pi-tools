package trips

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/database/migrate"
	"github.com/mjm/pi-tools/pkg/migrate/postgres"
)

func TestNewTracker(t *testing.T) {
	ctx := context.Background()
	clock := clockwork.NewFakeClock()

	t.Run("starts empty when there are no trips", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		tt, err := NewTracker(db, &fakeMessagesClient{})
		assert.NoError(t, err)

		assert.Nil(t, tt.currentTrip)
		assert.Empty(t, tt.lastLeft)
		assert.Empty(t, tt.lastReturned)
	})

	t.Run("populates last left and last returned when there are previous trips", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		q := database.New(db)

		leftAt := clock.Now()
		trip, err := q.BeginTrip(ctx, database.BeginTripParams{
			ID:     uuid.New(),
			LeftAt: leftAt,
		})
		assert.NoError(t, err)

		clock.Advance(5 * time.Minute)
		returnedAt := clock.Now()
		assert.NoError(t, q.EndTrip(ctx, database.EndTripParams{
			ID:         trip.ID,
			ReturnedAt: sql.NullTime{Time: returnedAt, Valid: true},
		}))

		tt, err := NewTracker(db, &fakeMessagesClient{})
		assert.NoError(t, err)

		assert.Nil(t, tt.currentTrip)
		assert.Equal(t, leftAt, tt.lastLeft.UTC())
		assert.Equal(t, returnedAt, tt.lastReturned.UTC())
	})

	t.Run("populates last left when there is only one trip and it's in-progress", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		q := database.New(db)

		leftAt := clock.Now()
		trip, err := q.BeginTrip(ctx, database.BeginTripParams{
			ID:     uuid.New(),
			LeftAt: leftAt,
		})
		assert.NoError(t, err)

		tt, err := NewTracker(db, &fakeMessagesClient{})
		assert.NoError(t, err)

		assert.NotNil(t, tt.currentTrip)
		assert.Equal(t, trip, *tt.currentTrip)
		assert.Equal(t, leftAt, tt.lastLeft.UTC())
		assert.Empty(t, tt.lastReturned.UTC())
	})

	t.Run("populates last left and last returned when there are previous trips and a current trip", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		q := database.New(db)

		leftAt := clock.Now()
		trip, err := q.BeginTrip(ctx, database.BeginTripParams{
			ID:     uuid.New(),
			LeftAt: leftAt,
		})
		assert.NoError(t, err)

		clock.Advance(5 * time.Minute)
		returnedAt := clock.Now()
		assert.NoError(t, q.EndTrip(ctx, database.EndTripParams{
			ID:         trip.ID,
			ReturnedAt: sql.NullTime{Time: returnedAt, Valid: true},
		}))

		clock.Advance(4 * time.Hour)
		leftAt = clock.Now()
		trip, err = q.BeginTrip(ctx, database.BeginTripParams{
			ID:     uuid.New(),
			LeftAt: leftAt,
		})
		assert.NoError(t, err)

		tt, err := NewTracker(db, &fakeMessagesClient{})
		assert.NoError(t, err)

		assert.NotNil(t, tt.currentTrip)
		assert.Equal(t, trip, *tt.currentTrip)
		assert.Equal(t, leftAt, tt.lastLeft.UTC())
		assert.Equal(t, returnedAt, tt.lastReturned.UTC())
	})
}

func TestTracker_OnLeave(t *testing.T) {
	ctx := context.Background()
	clock := clockwork.NewFakeClock()

	t.Run("starts a new trip", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))
		q := database.New(db)

		tt, err := NewTracker(db, &fakeMessagesClient{})
		assert.NoError(t, err)
		tt.clock = clock

		tt.OnLeave(ctx, nil)

		trip, err := q.GetCurrentTrip(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, tt.currentTrip)
		assert.Equal(t, trip, *tt.currentTrip)
	})

	t.Run("does not start a trip if one is already in-progress", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		q := database.New(db)

		leftAt := clock.Now()
		trip, err := q.BeginTrip(ctx, database.BeginTripParams{
			ID:     uuid.New(),
			LeftAt: leftAt,
		})
		assert.NoError(t, err)

		tt, err := NewTracker(db, &fakeMessagesClient{})
		assert.NoError(t, err)

		tt.OnLeave(ctx, nil)
		assert.Equal(t, trip, *tt.currentTrip)

		trips, err := q.ListTrips(ctx)
		assert.NoError(t, err)
		assert.Len(t, trips, 1)
	})
}
