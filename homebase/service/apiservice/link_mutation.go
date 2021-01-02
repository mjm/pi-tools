package apiservice

import (
	"context"

	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func (r *Resolver) CreateLink(ctx context.Context, args struct {
	Input struct {
		ShortURL       string
		DestinationURL string
		Description    string
	}
}) (
	resp struct {
		Link *Link
	},
	err error,
) {
	var res *linkspb.CreateLinkResponse
	res, err = r.linksClient.CreateLink(ctx, &linkspb.CreateLinkRequest{
		ShortUrl:       args.Input.ShortURL,
		DestinationUrl: args.Input.DestinationURL,
		Description:    args.Input.Description,
	})
	if err != nil {
		return
	}

	resp.Link = &Link{Link: res.GetLink()}
	return
}
