package backupservice

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
)

func (s *Server) GetArchive(ctx context.Context, req *backuppb.GetArchiveRequest) (*backuppb.GetArchiveResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("archive.id", req.GetId()),
		attribute.Stringer("archive.kind", req.GetKind()))

	resp := &backuppb.GetArchiveResponse{}

	switch req.GetKind() {
	case backuppb.Archive_BORG:
		a, err := s.borg.GetArchive(ctx, s.cfg.BorgRepoPath, req.GetId())
		if err != nil {
			return nil, err
		}

		startTs, err := ptypes.TimestampProto(a.StartTime.Time)
		if err != nil {
			return nil, err
		}
		endTs, err := ptypes.TimestampProto(a.EndTime.Time)
		if err != nil {
			return nil, err
		}

		resp.Archive = &backuppb.ArchiveDetail{
			Kind:        backuppb.Archive_BORG,
			Id:          a.ID,
			Name:        a.Name,
			StartTime:   startTs,
			EndTime:     endTs,
			Duration:    a.Duration,
			CommandLine: a.CommandLine,
			Username:    a.Username,
			Stats: &backuppb.ArchiveStats{
				CompressedSize:   a.Stats.CompressedSize,
				DeduplicatedSize: a.Stats.DeduplicatedSize,
				OriginalSize:     a.Stats.OriginalSize,
				NumFiles:         a.Stats.NumFiles,
			},
		}
		return resp, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "must specify a valid archive kind")
	}
}
