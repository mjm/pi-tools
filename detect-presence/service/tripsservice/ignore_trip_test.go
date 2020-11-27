package tripsservice

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mjm/pi-tools/detect-presence/database"
	"github.com/mjm/pi-tools/detect-presence/database/migrate"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/pkg/migrate/postgres"
)

func TestServer_IgnoreTrip(t *testing.T) {
	ctx := context.Background()

	t.Run("missing trip ID", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		// no need to migrate, we shouldn't make it to the query

		s := New(db, fakeMessagesClient{})
		res, err := s.IgnoreTrip(ctx, &tripspb.IgnoreTripRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing ID for trip to ignore")
	})

	t.Run("non-UUID trip ID", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		// no need to migrate, we shouldn't make it to the query

		s := New(db, fakeMessagesClient{})
		res, err := s.IgnoreTrip(ctx, &tripspb.IgnoreTripRequest{Id: "some nonsense"})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid UUID for trip ID: invalid UUID length: 13")
	})

	t.Run("missing trip", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		s := New(db, fakeMessagesClient{})
		id := uuid.New()
		res, err := s.IgnoreTrip(ctx, &tripspb.IgnoreTripRequest{Id: id.String()})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = no trip found with ID "+id.String())
	})

	t.Run("valid trip", func(t *testing.T) {
		db, err := dbSrv.NewDatabase(ctx)
		assert.NoError(t, err)
		assert.NoError(t, postgres.UpIfNeeded(db, migrate.Data))

		q := database.New(db)
		id := uuid.New()
		_, err = q.BeginTrip(ctx, database.BeginTripParams{
			ID:     id,
			LeftAt: time.Now(),
		})
		assert.NoError(t, err)

		s := New(db, fakeMessagesClient{})
		res, err := s.IgnoreTrip(ctx, &tripspb.IgnoreTripRequest{Id: id.String()})
		assert.NoError(t, err)
		assert.Equal(t, &tripspb.IgnoreTripResponse{}, res)

		// check that the trip was actually ignored by listing trips and it not being included
		trips, err := q.ListTrips(ctx, 30)
		assert.NoError(t, err)
		assert.Empty(t, trips)
	})
}
