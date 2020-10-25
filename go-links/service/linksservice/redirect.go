package linksservice

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/status"
)

func (s *Server) HandleShortLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	shortURL := strings.Trim(r.URL.Path, "/")
	span.SetAttributes(label.String("link.short_url", shortURL))

	link, err := s.db.GetLinkByShortURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("No go link found for go/%s", shortURL), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}

		s, _ := status.FromError(err)
		span.SetStatus(codes.Error, s.Message())
		return
	}

	span.SetAttributes(
		label.String("link.id", link.ID.String()),
		label.String("link.destination_url", link.DestinationURL))
	http.Redirect(w, r, link.DestinationURL, http.StatusPermanentRedirect)
}
