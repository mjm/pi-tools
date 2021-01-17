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
	githubRepo   = flag.String("repo", "mjm/pi-tools", "GitHub repository to check for builds")
	githubBranch = flag.String("branch", "main", "Git branch of builds that should be considered for deploying")
	workflowName = flag.String("workflow", "images.yaml", "Filename of GitHub Actions workflow to wait for")

	githubTokenPath = flag.String("github-token-path", "/secrets/github-token", "Path to file containing GitHub PAT token")

	terraformPath = flag.String("terraform", "/terraform", "Path to the Terraform binary")

	pollInterval = flag.Duration("poll-interval", 2*time.Minute, "How often to check with GitHub for a new build")
	dryRun       = flag.Bool("dry-run", false, "Skip actually applying changes to the cluster")
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
		DryRun:        *dryRun,
		GitHubRepo:    *githubRepo,
		GitHubBranch:  *githubBranch,
		WorkflowName:  *workflowName,
		TerraformPath: *terraformPath,
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
