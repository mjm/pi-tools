package apiservice

import (
	"net/http"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
)

type Server struct {
	handler http.Handler
}

func New(schemaString string, trips tripspb.TripsServiceClient, deploys deploypb.DeployServiceClient) (*Server, error) {
	r := &Resolver{
		tripsClient:  trips,
		deployClient: deploys,
	}
	schema, err := graphql.ParseSchema(schemaString, r, graphql.UseFieldResolvers())
	if err != nil {
		return nil, err
	}
	return &Server{
		handler: &relay.Handler{Schema: schema},
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
