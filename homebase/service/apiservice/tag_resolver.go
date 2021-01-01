package apiservice

import (
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

type TagConnection struct {
	res *tripspb.ListTagsResponse
}

func (tc *TagConnection) Edges() []TagEdge {
	var edges []TagEdge
	for _, t := range tc.res.GetTags() {
		edges = append(edges, TagEdge{Tag: t})
	}
	return edges
}

func (tc *TagConnection) PageInfo() PageInfo {
	return PageInfo{}
}

func (tc *TagConnection) TotalCount() int32 {
	return int32(len(tc.res.GetTags()))
}

type TagEdge struct {
	*tripspb.Tag
}

func (e TagEdge) Node() *Tag {
	return &Tag{Tag: e.Tag}
}

func (e TagEdge) Cursor() Cursor {
	return Cursor(e.GetName())
}

type Tag struct {
	*tripspb.Tag
}

func (t *Tag) Name() string {
	return t.GetName()
}

func (t *Tag) TripCount() int32 {
	return int32(t.GetTripCount())
}
