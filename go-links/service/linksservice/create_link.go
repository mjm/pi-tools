package linksservice

import (
	"context"

	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func (s *Server) CreateLink(ctx context.Context, req *linkspb.CreateLinkRequest) (*linkspb.CreateLinkResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("link.short_url", req.GetShortUrl()),
		attribute.String("link.destination_url", req.GetDestinationUrl()))

	if req.GetShortUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "short URL of link cannot be empty")
	}
	if req.GetDestinationUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "destination URL of link cannot be empty")
	}

	id := ksuid.New()
	span.SetAttributes(attribute.String("link.id", id.String()))

	link, err := s.db.CreateLink(ctx, database.CreateLinkParams{
		ID:             id,
		ShortURL:       req.GetShortUrl(),
		DestinationURL: req.GetDestinationUrl(),
		Description:    req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}
	linksCreatedTotal.Add(ctx, 1)

	return &linkspb.CreateLinkResponse{
		Link: marshalLinkToProto(link),
	}, nil
}
