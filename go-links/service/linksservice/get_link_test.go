package linksservice

import (
	"context"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"

	"github.com/mjm/pi-tools/go-links/database"
	"github.com/mjm/pi-tools/go-links/database/migrate"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
	"github.com/mjm/pi-tools/storage/storagetest"
)

func TestServer_GetLink(t *testing.T) {
	ctx := context.Background()

	t.Run("missing ID", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db)
		res, err := s.GetLink(ctx, &linkspb.GetLinkRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = missing ID of link to get")
	})

	t.Run("invalid ID", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db)
		res, err := s.GetLink(ctx, &linkspb.GetLinkRequest{Id: "garbage"})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = invalid link ID garbage: Valid encoded KSUIDs are 27 characters")
	})

	t.Run("missing link", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db)
		id := ksuid.New()
		res, err := s.GetLink(ctx, &linkspb.GetLinkRequest{Id: id.String()})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = no link found with ID "+id.String())
	})

	t.Run("valid link", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		q := database.New(db)
		id := ksuid.New()
		_, err = q.CreateLink(ctx, database.CreateLinkParams{
			ID:             id,
			ShortURL:       "foo",
			DestinationURL: "http://example.org",
			Description:    "A new link I just made",
		})
		assert.NoError(t, err)

		s := New(db)
		res, err := s.GetLink(ctx, &linkspb.GetLinkRequest{Id: id.String()})
		assert.NoError(t, err)
		assert.Equal(t, res, &linkspb.GetLinkResponse{
			Link: &linkspb.Link{
				Id:             id.String(),
				ShortUrl:       "foo",
				DestinationUrl: "http://example.org",
				Description:    "A new link I just made",
			},
		})
	})
}
