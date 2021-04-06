package apiservice

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"golang.org/x/sync/errgroup"

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

	return &Deploy{
		Deploy: res.GetDeploy(),
		r:      r,
	}, nil
}

func (r *Resolver) Deploy(ctx context.Context, args struct {
	ID string
}) (*Deploy, error) {
	if err := requireAuthorizedUser(ctx); err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(args.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	res, err := r.deployClient.GetDeploy(ctx, &deploypb.GetDeployRequest{
		DeployId: id,
	})
	if err != nil {
		return nil, err
	}

	return &Deploy{
		Deploy: res.GetDeploy(),
		r:      r,
	}, nil
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
	var borgArchives, tarsnapArchives []*backuppb.Archive

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		res, err := r.backupClient.ListArchives(groupCtx, &backuppb.ListArchivesRequest{
			Kind: backuppb.Archive_BORG,
		})
		if err != nil {
			return err
		}

		borgArchives = res.GetArchives()
		return nil
	})

	group.Go(func() error {
		res, err := r.backupClient.ListArchives(groupCtx, &backuppb.ListArchivesRequest{
			Kind: backuppb.Archive_TARSNAP,
		})
		if err != nil {
			return err
		}

		tarsnapArchives = res.GetArchives()
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	var archives []*backuppb.Archive
	var bi, ti int
	for bi < len(borgArchives) && ti < len(borgArchives) {
		if borgArchives[bi].GetTime().AsTime().After(tarsnapArchives[ti].GetTime().AsTime()) {
			archives = append(archives, borgArchives[bi])
			bi++
		} else {
			archives = append(archives, tarsnapArchives[ti])
			ti++
		}
	}

	for ; bi < len(borgArchives); bi++ {
		archives = append(archives, borgArchives[bi])
	}

	for ; ti < len(tarsnapArchives); ti++ {
		archives = append(archives, tarsnapArchives[ti])
	}

	return &ArchiveConnection{archives: archives}, nil
}

func (r *Resolver) BackupArchive(ctx context.Context, args struct{ ID graphql.ID }) (*ArchiveDetails, error) {
	switch relay.UnmarshalKind(args.ID) {
	case "borg_archive":
		var id string
		if err := relay.UnmarshalSpec(args.ID, &id); err != nil {
			return nil, err
		}

		res, err := r.backupClient.GetArchive(ctx, &backuppb.GetArchiveRequest{
			Kind: backuppb.Archive_BORG,
			Id:   id,
		})
		if err != nil {
			return nil, err
		}

		return &ArchiveDetails{ArchiveDetail: res.GetArchive()}, nil
	default:
		return nil, fmt.Errorf("unsupported kind")
	}
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
