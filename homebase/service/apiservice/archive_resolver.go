package apiservice

import (
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
	return relay.MarshalID("archive", a.GetId())
}

func (a *Archive) RawID() string {
	return a.GetId()
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
