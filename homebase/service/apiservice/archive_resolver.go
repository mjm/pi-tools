package apiservice

import (
	"fmt"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
)

type ArchiveConnection struct {
	archives []*backuppb.Archive
}

func (ac *ArchiveConnection) Edges() []ArchiveEdge {
	var edges []ArchiveEdge
	for _, a := range ac.archives {
		edges = append(edges, ArchiveEdge{Archive: a})
	}
	return edges
}

func (ArchiveConnection) PageInfo() PageInfo {
	return PageInfo{}
}

func (ac *ArchiveConnection) TotalCount() int32 {
	return int32(len(ac.archives))
}

type ArchiveEdge struct {
	*backuppb.Archive
}

func (e ArchiveEdge) Node() *Archive {
	return &Archive{Archive: e.Archive}
}

func (e ArchiveEdge) Cursor() Cursor {
	return Cursor(e.GetId())
}

type Archive struct {
	*backuppb.Archive
}

func (a *Archive) ID() graphql.ID {
	switch a.GetKind() {
	case backuppb.Archive_BORG:
		return relay.MarshalID("borg_archive", a.GetName())
	case backuppb.Archive_TARSNAP:
		return relay.MarshalID("tarsnap_archive", a.GetId())
	case backuppb.Archive_UNKNOWN:
		return relay.MarshalID("archive", a.GetId())
	}
	return ""
}

func (a *Archive) Kind() string {
	switch a.GetKind() {
	case backuppb.Archive_BORG:
		return "BORG"
	case backuppb.Archive_TARSNAP:
		return "TARSNAP"
	default:
		return "UNKNOWN"
	}
}

func (a *Archive) Name() string {
	return a.GetName()
}

func (a *Archive) CreatedAt() graphql.Time {
	return graphql.Time{Time: a.GetTime().AsTime()}
}

func (a *Archive) Details() (*ArchiveDetails, error) {
	return nil, fmt.Errorf("no details included in this archive")
}

func (a *Archive) Stats() (*ArchiveStats, error) {
	return nil, fmt.Errorf("no stats included in this archive")
}

type ArchiveDetails struct {
	*backuppb.ArchiveDetail
}

func (a *ArchiveDetails) ID() graphql.ID {
	switch a.GetKind() {
	case backuppb.Archive_BORG:
		return relay.MarshalID("borg_archive", a.GetName())
	case backuppb.Archive_TARSNAP:
		return relay.MarshalID("tarsnap_archive", a.GetId())
	case backuppb.Archive_UNKNOWN:
		return relay.MarshalID("archive", a.GetId())
	}
	return ""
}

func (a *ArchiveDetails) Kind() string {
	switch a.GetKind() {
	case backuppb.Archive_BORG:
		return "BORG"
	case backuppb.Archive_TARSNAP:
		return "TARSNAP"
	default:
		return "UNKNOWN"
	}
}

func (a *ArchiveDetails) Name() string {
	return a.GetName()
}

func (a *ArchiveDetails) CreatedAt() graphql.Time {
	return graphql.Time{Time: a.GetStartTime().AsTime()}
}

func (a *ArchiveDetails) Details() *ArchiveDetails {
	return a
}

func (a *ArchiveDetails) Stats() *ArchiveStats {
	return &ArchiveStats{ArchiveStats: a.GetStats()}
}

func (a *ArchiveDetails) Duration() float64 {
	return a.GetDuration()
}

func (a *ArchiveDetails) CommandLine() []string {
	return a.GetCommandLine()
}

type ArchiveStats struct {
	*backuppb.ArchiveStats
}

func (a *ArchiveStats) CompressedSize() int32 {
	return int32(a.GetCompressedSize())
}

func (a *ArchiveStats) DeduplicatedSize() int32 {
	return int32(a.GetDeduplicatedSize())
}

func (a *ArchiveStats) OriginalSize() int32 {
	return int32(a.GetOriginalSize())
}

func (a *ArchiveStats) NumFiles() int32 {
	return int32(a.GetNumFiles())
}
