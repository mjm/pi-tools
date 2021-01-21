package apiservice

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	"github.com/mjm/pi-tools/debug"
	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	linkspb "github.com/mjm/pi-tools/go-links/proto/links"
)

type Server struct {
	handler http.Handler
}

func New(
	schemaString string,
	trips tripspb.TripsServiceClient,
	links linkspb.LinksServiceClient,
	deploys deploypb.DeployServiceClient,
	prometheusURL string,
) (*Server, error) {
	r := &Resolver{
		tripsClient:   trips,
		linksClient:   links,
		deployClient:  deploys,
		prometheusURL: prometheusURL,
	}
	schema, err := graphql.ParseSchema(schemaString, r, graphql.UseFieldResolvers())
	if err != nil {
		return nil, err
	}
	return &Server{
		handler: &relay.Handler{Schema: schema},
	}, nil
}

type contextKey string

const (
	cookieHeaderContextKey contextKey = "cookie-header"
	authUserContextKey     contextKey = "auth-user"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), cookieHeaderContextKey, r.Header.Get("Cookie"))
	ctx = context.WithValue(ctx, authUserContextKey, r.Header.Get("X-Auth-Request-User"))
	r = r.WithContext(ctx)
	s.handler.ServeHTTP(w, r)
}

func requireAuthorizedUser(ctx context.Context) error {
	if debug.IsEnabled() {
		return nil
	}

	user := ctx.Value(authUserContextKey).(string)
	if user == "" {
		return fmt.Errorf("this data is only available to authorized users")
	}
	return nil
}
