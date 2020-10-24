package linksservice

import (
	"context"

	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func (s *Server) CreateLink(ctx context.Context, req *linkspb.CreateLinkRequest) (*linkspb.CreateLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "CreateLink",
		trace.WithAttributes(
			label.String("link.short_url", req.GetShortUrl()),
			label.String("link.destination_url", req.GetDestinationUrl())))
	defer span.End()

	id := ksuid.New()
	span.SetAttributes(label.String("link.id", id.String()))

	link, err := s.db.CreateLink(ctx, database.CreateLinkParams{
		ID:             id,
		ShortURL:       req.GetShortUrl(),
		DestinationURL: req.GetDestinationUrl(),
		Description:    req.GetDescription(),
	})
	if err != nil {
		span.RecordError(ctx, err)
		return nil, err
	}
	linksCreatedTotal.Add(ctx, 1)

	return &linkspb.CreateLinkResponse{
		Link: marshalLinkToProto(link),
	}, nil
}
