package apiservice

import (
	"context"
	"net/http"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	backuppb "github.com/mjm/pi-tools/backup/proto/backup"
	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

type Resolver struct {
	tripsClient   tripspb.TripsServiceClient
	linksClient   linkspb.LinksServiceClient
	deployClient  deploypb.DeployServiceClient
	backupClient  backuppb.BackupServiceClient
	prometheusURL string
}

func (r *Resolver) Viewer() *Resolver {
	return r
}

func (r *Resolver) Node(ctx context.Context, args struct{ ID graphql.ID }) (*Node, error) {
	return nil, nil
}

func (r *Resolver) Trips(ctx context.Context, args struct {
	First *int32
	After *Cursor
}) (*TripConnection, error) {
	// TODO actually support paging

	var limit int32 = 30
	if args.First != nil {
		limit = *args.First
	}
	res, err := r.tripsClient.ListTrips(ctx, &tripspb.ListTripsRequest{Limit: limit})
	if err != nil {
		return nil, err
	}
	return &TripConnection{res: res}, nil
}

func (r *Resolver) Trip(ctx context.Context, args struct {
	ID graphql.ID
}) (*Trip, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	var id string
	if err := relay.UnmarshalSpec(args.ID, &id); err != nil {
		return nil, err
	}

	res, err := r.tripsClient.GetTrip(ctx, &tripspb.GetTripRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &Trip{Trip: res.GetTrip()}, nil
}

func (r *Resolver) Tags(ctx context.Context, args struct {
	First *int32
	After *Cursor
}) (*TagConnection, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	// TODO actually support paging

	var limit int32 = 30
	if args.First != nil {
		limit = *args.First
	}
	res, err := r.tripsClient.ListTags(ctx, &tripspb.ListTagsRequest{Limit: limit})
	if err != nil {
		return nil, err
	}
	return &TagConnection{res: res}, nil
}

func (r *Resolver) Links(ctx context.Context, args struct {
	First *int32
	After *Cursor
}) (*LinkConnection, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	// TODO actually support paging

	//var limit int32 = 30
	//if args.First != nil {
	//	limit = *args.First
	//}
	res, err := r.linksClient.ListRecentLinks(ctx, &linkspb.ListRecentLinksRequest{})
	if err != nil {
		return nil, err
	}
	return &LinkConnection{res: res}, nil
}

func (r *Resolver) Link(ctx context.Context, args struct {
	ID graphql.ID
}) (*Link, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	var id string
	if err := relay.UnmarshalSpec(args.ID, &id); err != nil {
		return nil, err
	}

	res, err := r.linksClient.GetLink(ctx, &linkspb.GetLinkRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &Link{Link: res.GetLink()}, nil
}

func (r *Resolver) MostRecentDeploy(ctx context.Context) (*Deploy, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	res, err := r.deployClient.GetMostRecentDeploy(ctx, &deploypb.GetMostRecentDeployRequest{})
	if err != nil {
		return nil, err
	}

	return &Deploy{Deploy: res.GetDeploy()}, nil
}

func (r *Resolver) Alerts(ctx context.Context) ([]*Alert, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	promClient, err := r.newPromClient(ctx)
	if err != nil {
		return nil, err
	}

	alerts, err := promClient.Alerts(ctx)
	if err != nil {
		return nil, err
	}

	var res []*Alert
	for _, a := range alerts.Alerts {
		res = append(res, &Alert{Alert: a})
	}
	return res, nil
}

func (r *Resolver) BackupArchives(ctx context.Context, args struct {
	First *int32
	After *Cursor
	Kind  *string
}) (*ArchiveConnection, error) {
	res, err := r.backupClient.ListArchives(ctx, &backuppb.ListArchivesRequest{})
	if err != nil {
		return nil, err
	}

	return &ArchiveConnection{res: res}, nil
}

func (r *Resolver) newPromClient(ctx context.Context) (v1.API, error) {
	transport := http.DefaultTransport
	if strings.HasPrefix(r.prometheusURL, "https://") {
		cookieHeader := ctx.Value(cookieHeaderContextKey).(string)
		transport = &oauthProxyCookieTripper{
			cookieHeader: cookieHeader,
			wrapped:      transport,
		}
	}

	c, err := api.NewClient(api.Config{
		Address:      r.prometheusURL,
		RoundTripper: transport,
	})
	if err != nil {
		return nil, err
	}
	return v1.NewAPI(c), nil
}

type oauthProxyCookieTripper struct {
	cookieHeader string
	wrapped      http.RoundTripper
}

func (t *oauthProxyCookieTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	transport := t.wrapped
	if transport == nil {
		transport = http.DefaultTransport
	}

	r.Header.Set("Cookie", t.cookieHeader)
	return transport.RoundTrip(r)
}
