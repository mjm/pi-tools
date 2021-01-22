package appservice

import (
	"net/http"

	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	owner = "mjm"
	repo  = "pi-tools"
)

type Server struct {
	HTTPClient   *http.Client
	GithubClient *github.Client
}

func New(githubClient *github.Client) *Server {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return &Server{
		HTTPClient:   httpClient,
		GithubClient: githubClient,
	}
}
