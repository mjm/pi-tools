package linksservice

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func (s *Server) ListRecentLinks(ctx context.Context, req *linkspb.ListRecentLinksRequest) (*linkspb.ListRecentLinksResponse, error) {
	span := trace.SpanFromContext(ctx)

	links, err := s.db.ListRecentLinks(ctx)
	if err != nil {
		return nil, err
	}
	span.SetAttributes(label.Int("link.count", len(links)))

	var res linkspb.ListRecentLinksResponse
	for _, link := range links {
		res.Links = append(res.Links, marshalLinkToProto(link))
	}
	return &res, nil
}
