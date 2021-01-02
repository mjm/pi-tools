package apiservice

import (
	"context"
	"net/http"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

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

const cookieHeaderContextKey contextKey = "cookie-header"

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), cookieHeaderContextKey, r.Header.Get("Cookie"))
	r = r.WithContext(ctx)
	s.handler.ServeHTTP(w, r)
}
