package linksservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func (s *Server) UpdateLink(ctx context.Context, req *linkspb.UpdateLinkRequest) (*linkspb.UpdateLinkResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		label.String("link.id", req.GetId()),
		label.String("link.short_url", req.GetShortUrl()),
		label.String("link.destination_url", req.GetDestinationUrl()))

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID of link to update")
	}
	if req.GetShortUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "short URL of link cannot be empty")
	}
	if req.GetDestinationUrl() == "" {
		return nil, status.Error(codes.InvalidArgument, "destination URL of link cannot be empty")
	}

	id, err := ksuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID %s: %v", req.GetId(), err)
	}

	link, err := s.db.UpdateLink(ctx, database.UpdateLinkParams{
		ID:             id,
		ShortURL:       req.GetShortUrl(),
		DestinationURL: req.GetDestinationUrl(),
		Description:    req.GetDescription(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no link found with ID %s", id)
		}
		return nil, err
	}

	return &linkspb.UpdateLinkResponse{
		Link: marshalLinkToProto(link),
	}, nil
}
