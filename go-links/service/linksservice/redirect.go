package linksservice

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) HandleShortLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shortURL := strings.Trim(r.URL.Path, "/")
	link, err := s.db.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("No go link found for go/%s", shortURL)))
		} else {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}

	http.Redirect(w, r, link.DestinationURL, http.StatusPermanentRedirect)
}
