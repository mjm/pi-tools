package borgbackup

import (
	"context"
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	t.Time, err = time.ParseInLocation("2006-01-02T15:04:05.999999", s, time.Local)
	if err != nil {
		return err
	}

	return nil
}

type ArchiveSummary struct {
	Archive string `json:"archive"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Time    Time   `json:"start"`
}

type listArchivesResponse struct {
	Archives []*ArchiveSummary `json:"archives"`
}

func (b *Borg) ListArchives(ctx context.Context, repositoryPath string) ([]*ArchiveSummary, error) {
	var resp listArchivesResponse
	if err := b.commandJSON(ctx, &resp, "list", repositoryPath, "--last", "10"); err != nil {
		return nil, err
	}

	return resp.Archives, nil
}
