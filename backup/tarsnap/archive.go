package tarsnap

import (
	"context"
	"sort"
	"strings"
	"time"
)

type ArchiveSummary struct {
	Name      string
	CreatedAt time.Time
}

func (t *Tarsnap) ListArchives(ctx context.Context, keyPath string) ([]*ArchiveSummary, error) {
	out, err := t.runCommand(ctx, "--keyfile", keyPath, "--list-archives", "-v")
	if err != nil {
		return nil, err
	}

	var summaries []*ArchiveSummary

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		tt, err := time.ParseInLocation("2006-01-02 15:04:05", fields[1], time.Local)
		if err != nil {
			return nil, err
		}

		summaries = append(summaries, &ArchiveSummary{
			Name:      fields[0],
			CreatedAt: tt,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[j].Name < summaries[i].Name
	})
	return summaries, nil
}
