package borgbackup

import (
	"context"
	"encoding/json"
	"fmt"
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

type Archive struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	StartTime   Time         `json:"start"`
	EndTime     Time         `json:"end"`
	Duration    float64      `json:"duration"`
	CommandLine []string     `json:"command_line"`
	Username    string       `json:"username"`
	Stats       ArchiveStats `json:"stats"`
}

type ArchiveStats struct {
	CompressedSize   int64 `json:"compressed_size"`
	DeduplicatedSize int64 `json:"deduplicated_size"`
	NumFiles         int64 `json:"nfiles"`
	OriginalSize     int64 `json:"original_size"`
}

type getArchiveResponse struct {
	Archives []*Archive `json:"archives"`
}

func (b *Borg) GetArchive(ctx context.Context, repositoryPath string, name string) (*Archive, error) {
	archivePath := fmt.Sprintf("%s::%s", repositoryPath, name)

	var resp getArchiveResponse
	if err := b.commandJSON(ctx, &resp, "info", archivePath); err != nil {
		return nil, err
	}

	return resp.Archives[0], nil
}
