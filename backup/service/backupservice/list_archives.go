package backupservice

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
)

func (s *Server) ListArchives(ctx context.Context, req *backuppb.ListArchivesRequest) (*backuppb.ListArchivesResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.Stringer("archive.kind", req.GetKind()))

	resp := &backuppb.ListArchivesResponse{}

	switch req.GetKind() {
	case backuppb.Archive_BORG:
		archives, err := s.borg.ListArchives(ctx, s.cfg.BorgRepoPath)
		if err != nil {
			return nil, err
		}

		span.SetAttributes(label.Int("archive.count", len(archives)))

		for i := len(archives) - 1; i >= 0; i-- {
			archive := archives[i]
			ts, err := ptypes.TimestampProto(archive.Time.Time)
			if err != nil {
				return nil, err
			}

			resp.Archives = append(resp.Archives, &backuppb.Archive{
				Kind: backuppb.Archive_BORG,
				Id:   archive.ID,
				Name: archive.Name,
				Time: ts,
			})
		}
		return resp, nil
	case backuppb.Archive_TARSNAP:
		archives, err := s.tarsnap.ListArchives(ctx, s.cfg.TarsnapKeyPath)
		if err != nil {
			return nil, err
		}

		span.SetAttributes(label.Int("archive.count", len(archives)))

		for _, archive := range archives {
			ts, err := ptypes.TimestampProto(archive.CreatedAt)
			if err != nil {
				return nil, err
			}

			resp.Archives = append(resp.Archives, &backuppb.Archive{
				Kind: backuppb.Archive_TARSNAP,
				Id:   archive.Name,
				Name: archive.Name,
				Time: ts,
			})
		}
		return resp, nil
	}

	return nil, status.Errorf(codes.InvalidArgument, "must specify a valid archive kind")
}
