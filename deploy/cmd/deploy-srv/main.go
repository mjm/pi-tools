package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/etherlabsio/healthcheck"
	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/deploy/service/deployservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	githubRepo   = flag.String("repo", "mjm/pi-tools", "GitHub repository to pull down artifact from")
	githubBranch = flag.String("branch", "main", "Git branch of builds that should be considered for deploying")
	workflowName = flag.String("workflow", "k8s.yaml", "Filename of GitHub Actions workflow to pull artifact from")
	artifactName = flag.String("artifact", "all_k8s", "Name of artifact to download from workflow run")
	fileToApply  = flag.String("f", "k8s", "Extension-less name of the YAML file to apply")

	githubTokenPath = flag.String("github-token-path", "/var/secrets/github-token", "Path to file containing GitHub PAT token")

	pollInterval = flag.Duration("poll-interval", 5*time.Minute, "How often to check with GitHub for a new build artifact")
	dryRun       = flag.Bool("dry-run", false, "Skip actually applying changes to the Kubernetes cluster")
)

func main() {
	rpc.SetDefaultHTTPPort(8480)
	rpc.SetDefaultGRPCPort(8481)
	flag.Parse()

	stopObs := observability.MustStart("deploy-srv")
	defer stopObs()

	tokenData, err := ioutil.ReadFile(*githubTokenPath)
	if err != nil {
		log.Panicf("reading github token: %v", err)
	}

	token := &oauth2.Token{AccessToken: strings.TrimSpace(string(tokenData))}
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)
	githubClient := github.NewClient(httpClient)

	deploysSrv := deployservice.New(githubClient, deployservice.Config{
		DryRun:       *dryRun,
		GitHubRepo:   *githubRepo,
		GitHubBranch: *githubBranch,
		WorkflowName: *workflowName,
		ArtifactName: *artifactName,
		FileToApply:  *fileToApply,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go deploysSrv.PollForChanges(ctx, *pollInterval)

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		deploypb.RegisterDeployServiceServer(s, deploysSrv)
	}))

	signal.Wait()
}
