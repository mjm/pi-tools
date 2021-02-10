package backupservice

import (
	"context"

	"github.com/golang/protobuf/ptypes"

	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
)

func (s *Server) ListArchives(ctx context.Context, req *backuppb.ListArchivesRequest) (*backuppb.ListArchivesResponse, error) {
	archives, err := s.borg.ListArchives(ctx, s.cfg.BorgRepoPath)
	if err != nil {
		return nil, err
	}

	resp := &backuppb.ListArchivesResponse{}
	for _, archive := range archives {
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
}
