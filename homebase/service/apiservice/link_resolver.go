package apiservice

import (
	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

type LinkConnection struct {
	res *linkspb.ListRecentLinksResponse
}

func (lc *LinkConnection) Edges() []LinkEdge {
	var edges []LinkEdge
	for _, l := range lc.res.GetLinks() {
		edges = append(edges, LinkEdge{Link: l})
	}
	return edges
}

func (LinkConnection) PageInfo() PageInfo {
	return PageInfo{}
}

func (lc *LinkConnection) TotalCount() int32 {
	return int32(len(lc.res.GetLinks()))
}

type LinkEdge struct {
	*linkspb.Link
}

func (e LinkEdge) Node() *Link {
	return &Link{Link: e.Link}
}

func (e LinkEdge) Cursor() Cursor {
	return Cursor(e.GetId())
}

type Link struct {
	*linkspb.Link
}

func (l *Link) ID() graphql.ID {
	return relay.MarshalID("link", l.GetId())
}

func (l *Link) RawID() string {
	return l.GetId()
}

func (l *Link) ShortURL() string {
	return l.GetShortUrl()
}

func (l *Link) DestinationURL() string {
	return l.GetDestinationUrl()
}

func (l *Link) Description() string {
	return l.GetDescription()
}
