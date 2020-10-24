package linksservice

import (
	"github.com/mjm/pi-tools/go-links/database"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

func marshalLinkToProto(l *database.Link) *linkspb.Link {
	return &linkspb.Link{
		Id:             l.ID,
		ShortUrl:       l.ShortURL,
		DestinationUrl: l.DestinationURL,
		Description:    l.Description,
	}
}
