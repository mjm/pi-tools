package linksservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func (s *Server) GetLink(ctx context.Context, req *linkspb.GetLinkRequest) (*linkspb.GetLinkResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.String("link.id", req.GetId()))

	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing ID of link to get")
	}

	id, err := ksuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid link ID %s: %v", req.GetId(), err)
	}

	link, err := s.db.GetLink(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "no link found with ID %s", id)
		}
		return nil, err
	}

	return &linkspb.GetLinkResponse{
		Link: marshalLinkToProto(link),
	}, nil
}
