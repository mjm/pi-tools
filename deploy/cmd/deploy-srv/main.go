package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/etherlabsio/healthcheck"
	"github.com/google/go-github/v33/github"
	"github.com/gregdel/pushover"
	nomadapi "github.com/hashicorp/nomad/api"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/deploy/service/deployservice"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/leader"
	"github.com/mjm/pi-tools/pkg/signal"
	"github.com/mjm/pi-tools/rpc"
)

var (
	githubRepo   = flag.String("repo", "mjm/pi-tools", "GitHub repository to check for builds")
	githubBranch = flag.String("branch", "main", "Git branch of builds that should be considered for deploying")
	workflowName = flag.String("workflow", "nomad.yaml", "Filename of GitHub Actions workflow to wait for")
	artifactName = flag.String("artifact", "nomadic", "Name of artifact containing Nomad job JSON files")

	githubTokenPath = flag.String("github-token-path", "/secrets/github-token", "Path to file containing GitHub PAT token")

	pollInterval = flag.Duration("poll-interval", 2*time.Minute, "How often to check with GitHub for a new build")
	dryRun       = flag.Bool("dry-run", false, "Skip actually applying changes to the cluster")

	minioURL     = flag.String("minio-url", "http://localhost:9000", "URL for accessing Minio for storing deploy reports")
	reportBucket = flag.String("report-bucket", "deploy-reports", "Bucket to use to store deploy reports")
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

	nomadClient, err := nomadapi.NewClient(nomadapi.DefaultConfig())
	if err != nil {
		log.Panicf("creating nomad client: %v", err)
	}

	pushoverClient := pushover.New(os.Getenv("PUSHOVER_TOKEN"))
	pushoverRecipient := pushover.NewRecipient(os.Getenv("PUSHOVER_USER_KEY"))

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         minioURL,
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(strings.HasPrefix(*minioURL, "http://")),
		S3ForcePathStyle: aws.Bool(true),
	}))
	s3Client := s3.New(sess)

	deploysSrv := deployservice.New(githubClient, nomadClient, pushoverClient, s3Client, deployservice.Config{
		DryRun:            *dryRun,
		GitHubRepo:        *githubRepo,
		GitHubBranch:      *githubBranch,
		WorkflowName:      *workflowName,
		ArtifactName:      *artifactName,
		PushoverRecipient: pushoverRecipient,
		ReportBucket:      *reportBucket,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	election, err := leader.NewElection(leader.Config{
		Key: "service/deploy/leader",
		OnAcquireLeader: func() {
			deploysSrv.PollForChanges(ctx, *pollInterval)
		},
	})
	if err != nil {
		log.Panicf("Error creating leader election: %v", err)
	}

	go election.Run(ctx)
	defer election.Stop()

	// Ensure the current deploy, if any, is complete before shutting down and giving up leadership
	defer deploysSrv.Shutdown()

	http.Handle("/healthz",
		otelhttp.WithRouteTag("CheckHealth", healthcheck.Handler(
			healthcheck.WithTimeout(3*time.Second))))

	go rpc.ListenAndServe(rpc.WithRegisteredServices(func(s *grpc.Server) {
		deploypb.RegisterDeployServiceServer(s, deploysSrv)
	}))

	signal.Wait()
}
