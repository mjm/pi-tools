package tripsservice

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
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/storage/storagetest"
)

func TestServer_ListTrips(t *testing.T) {
	ctx := context.Background()

	t.Run("empty list of trips", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db, fakeMessagesClient{})
		res, err := s.ListTrips(ctx, &tripspb.ListTripsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Empty(t, res.GetTrips())
	})

	t.Run("with a few trips", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		clock := clockwork.NewFakeClockAt(time.Date(2020, 11, 3, 0, 0, 0, 0, time.UTC))
		q := database.New(db)

		tripID1 := uuid.New()
		trip, err := q.BeginTrip(ctx, database.BeginTripParams{
			ID:     tripID1,
			LeftAt: clock.Now(),
		})
		assert.NoError(t, err)
		clock.Advance(10 * time.Minute)
		assert.NoError(t, q.EndTrip(ctx, database.EndTripParams{
			ID:         trip.ID,
			ReturnedAt: sql.NullTime{Time: clock.Now(), Valid: true},
		}))
		clock.Advance(4 * time.Hour)

		tripID2 := uuid.New()
		trip, err = q.BeginTrip(ctx, database.BeginTripParams{
			ID:     tripID2,
			LeftAt: clock.Now(),
		})
		assert.NoError(t, err)
		clock.Advance(15 * time.Minute)
		assert.NoError(t, q.EndTrip(ctx, database.EndTripParams{
			ID:         trip.ID,
			ReturnedAt: sql.NullTime{Time: clock.Now(), Valid: true},
		}))

		assert.NoError(t, q.UpdateTripTags(ctx, database.UpdateTripTagsParams{
			TripID:    tripID2,
			TagsToAdd: []string{"long trip", "dog walk"},
		}))

		s := New(db, fakeMessagesClient{})
		res, err := s.ListTrips(ctx, &tripspb.ListTripsRequest{})
		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, &tripspb.ListTripsResponse{
			Trips: []*tripspb.Trip{
				{
					Id:         tripID2.String(),
					LeftAt:     "2020-11-03T04:10:00Z",
					ReturnedAt: "2020-11-03T04:25:00Z",
					Tags:       []string{"dog walk", "long trip"},
				},
				{
					Id:         tripID1.String(),
					LeftAt:     "2020-11-03T00:00:00Z",
					ReturnedAt: "2020-11-03T00:10:00Z",
					Tags:       []string{},
				},
			},
		}, res)
	})
}
