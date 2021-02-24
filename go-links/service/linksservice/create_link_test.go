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

func TestServer_CreateLink(t *testing.T) {
	ctx := context.Background()

	t.Run("missing short URL", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db)
		res, err := s.CreateLink(ctx, &linkspb.CreateLinkRequest{
			DestinationUrl: "http://example.org",
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = short URL of link cannot be empty")
	})

	t.Run("missing destination URL", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db)
		res, err := s.CreateLink(ctx, &linkspb.CreateLinkRequest{
			ShortUrl: "foo",
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = destination URL of link cannot be empty")
	})

	t.Run("valid input", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.FS)
		assert.NoError(t, err)

		s := New(db)
		res, err := s.CreateLink(ctx, &linkspb.CreateLinkRequest{
			ShortUrl:       "foo",
			DestinationUrl: "http://example.org",
			Description:    "A new link I just made",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.GetLink().GetId())
		assert.Equal(t, "foo", res.GetLink().GetShortUrl())
		assert.Equal(t, "http://example.org", res.GetLink().GetDestinationUrl())
		assert.Equal(t, "A new link I just made", res.GetLink().GetDescription())

		id, err := ksuid.Parse(res.GetLink().GetId())
		assert.NoError(t, err)
		link, err := database.New(db).GetLink(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, "foo", link.ShortURL)
	})
}
