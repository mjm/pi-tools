package linksservice

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"

	"github.com/mjm/pi-tools/go-links/database"
	"github.com/mjm/pi-tools/go-links/database/migrate"
	"github.com/mjm/pi-tools/storage/storagetest"
)

func TestServer_HandleShortLink(t *testing.T) {
	ctx := context.Background()

	t.Run("missing link", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.Data)
		assert.NoError(t, err)

		s := New(db)
		ts := httptest.NewServer(http.HandlerFunc(s.HandleShortLink))
		defer ts.Close()

		res, err := http.Get(ts.URL + "/foo")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		assert.NoError(t, err)
		assert.Equal(t, string(body), "No go link found for go/foo\n")
	})

	t.Run("valid link", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.Data)
		assert.NoError(t, err)

		q := database.New(db)
		_, err = q.CreateLink(ctx, database.CreateLinkParams{
			ID:             ksuid.New(),
			ShortURL:       "foo",
			DestinationURL: "http://example.org/foo",
		})
		assert.NoError(t, err)

		s := New(db)
		ts := httptest.NewServer(http.HandlerFunc(s.HandleShortLink))
		defer ts.Close()

		h := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		res, err := h.Get(ts.URL + "/foo")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusPermanentRedirect, res.StatusCode)
		assert.Equal(t, "http://example.org/foo", res.Header.Get("Location"))
	})

	t.Run("strips slashes from the short URL", func(t *testing.T) {
		db, err := storagetest.NewDatabase(ctx, dbSrv, migrate.Data)
		assert.NoError(t, err)

		q := database.New(db)
		_, err = q.CreateLink(ctx, database.CreateLinkParams{
			ID:             ksuid.New(),
			ShortURL:       "foo",
			DestinationURL: "http://example.org/foo",
		})
		assert.NoError(t, err)

		s := New(db)
		ts := httptest.NewServer(http.HandlerFunc(s.HandleShortLink))
		defer ts.Close()

		h := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		res, err := h.Get(ts.URL + "/foo/")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusPermanentRedirect, res.StatusCode)
		assert.Equal(t, "http://example.org/foo", res.Header.Get("Location"))
	})
}
