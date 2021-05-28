package apiservice

import (
	"encoding/base64"
	"net/url"
	"strconv"
	"time"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"
)

type PaperlessDocumentConnection struct {
	res *paperlessDocumentsResponse
	r   *Resolver
}

func (c *PaperlessDocumentConnection) Edges() []PaperlessDocumentEdge {
	var edges []PaperlessDocumentEdge
	for _, d := range c.res.Results {
		d := d
		edges = append(edges, PaperlessDocumentEdge{
			res: &d,
			r:   c.r,
		})
	}
	return edges
}

func (c *PaperlessDocumentConnection) PageInfo() (*PageInfo, error) {
	if c.res.Next == nil {
		return &PageInfo{}, nil
	}

	nextURL, err := url.Parse(*c.res.Next)
	if err != nil {
		return nil, err
	}
	q := nextURL.Query()

	encodedPage := base64.StdEncoding.EncodeToString([]byte(q.Get("page")))
	endCursor := Cursor(encodedPage)
	return &PageInfo{
		HasNextPage: false,
		EndCursor:   &endCursor,
	}, nil
}

func (c *PaperlessDocumentConnection) TotalCount() int32 {
	return int32(c.res.Count)
}

type PaperlessDocumentEdge struct {
	res *paperlessDocumentResponse
	r   *Resolver
}

func (e PaperlessDocumentEdge) Node() *PaperlessDocument {
	return &PaperlessDocument{
		res: e.res,
		r:   e.r,
	}
}

func (e PaperlessDocumentEdge) Cursor() Cursor {
	return Cursor(strconv.Itoa(e.res.ID))
}

type PaperlessDocument struct {
	res *paperlessDocumentResponse
	r   *Resolver
}

func (d *PaperlessDocument) ID() graphql.ID {
	return relay.MarshalID("paperless_doc", d.res.ID)
}

func (d *PaperlessDocument) Title() string {
	return d.res.Title
}

func (d *PaperlessDocument) CreatedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, d.res.Created)
	if err != nil {
		return graphql.Time{}, err
	}
	return graphql.Time{Time: t}, nil
}

func (d *PaperlessDocument) AddedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, d.res.Added)
	if err != nil {
		return graphql.Time{}, err
	}
	return graphql.Time{Time: t}, nil
}

func (d *PaperlessDocument) ModifiedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, d.res.Modified)
	if err != nil {
		return graphql.Time{}, err
	}
	return graphql.Time{Time: t}, nil
}
