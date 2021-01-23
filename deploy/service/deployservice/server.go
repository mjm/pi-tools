package deployservice

import (
	"github.com/google/go-github/v33/github"
	nomadapi "github.com/hashicorp/nomad/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// Config contains configuration parameters for how to fetch the build artifact to deploy.
type Config struct {
	// DryRun is a flag controlling whether to actually attempt to update the Kubernetes cluster.
	// If it is true, all actions will be performed up until the actual kubectl command would be
	// run, and then that will be skipped and assumed to have succeeded.
	DryRun bool

	// GitHubRepo is the full name of the GitHub repository to pull the build artifact from.
	GitHubRepo string

	// GitHubBranch is the branch whose builds should be considered for deploying.
	GitHubBranch string

	// WorkflowName is the filename of the GitHub Actions workflow whose artifact should be used.
	WorkflowName string

	// ArtifactName is the name of the artifact in the workflow that contains the JSON files for
	// the Nomad jobs to deploy.
	ArtifactName string
}

type Server struct {
	// Config is the configuration for deployment.
	Config Config

	// GitHubClient is the client to use to make API requests to GitHub.
	GitHubClient *github.Client

	// NomadClient is the client to use to submit jobs to Nomad
	NomadClient *nomadapi.Client

	lastSuccessfulCommit string
	deployChecksTotal    metric.Int64Counter
	deployCheckDuration  metric.Float64ValueRecorder
}

func New(gh *github.Client, nomad *nomadapi.Client, cfg Config) *Server {
	m := metric.Must(otel.Meter(instrumentationName))
	return &Server{
		Config:       cfg,
		GitHubClient: gh,
		NomadClient:  nomad,
		deployChecksTotal: m.NewInt64Counter("deploy.check.total",
			metric.WithDescription("Counts the number of times that the service checked for a new version to deploy")),
		deployCheckDuration: m.NewFloat64ValueRecorder("deploy.check.duration.seconds",
			metric.WithDescription("Records the amount of time spent checking for and applying new changes")),
	}
}
