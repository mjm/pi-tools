package apiservice

import (
	"context"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

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
	if err = requireAuthorizedUser(ctx); err != nil {
		return
	}

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

func (r *Resolver) UpdateLink(ctx context.Context, args struct {
	Input struct {
		ID             graphql.ID
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
	if err = requireAuthorizedUser(ctx); err != nil {
		return
	}

	var id string
	if err = relay.UnmarshalSpec(args.Input.ID, &id); err != nil {
		return
	}

	var res *linkspb.UpdateLinkResponse
	res, err = r.linksClient.UpdateLink(ctx, &linkspb.UpdateLinkRequest{
		Id:             id,
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
