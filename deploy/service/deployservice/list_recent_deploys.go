package deployservice

import (
	"context"
	"strings"

	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

func (s *Server) ListRecentDeploys(ctx context.Context, _ *deploypb.ListRecentDeploysRequest) (*deploypb.ListRecentDeploysResponse, error) {
	span := trace.SpanFromContext(ctx)

	repoParts := strings.SplitN(s.Config.GitHubRepo, "/", 2)
	deployments, _, err := s.GitHubClient.Repositories.ListDeployments(ctx, repoParts[0], repoParts[1], &github.DeploymentsListOptions{
		Task:        "deploy",
		Environment: "production",
		ListOptions: github.ListOptions{
			PerPage: 10,
		},
	})
	if err != nil {
		return nil, err
	}

	span.SetAttributes(label.Int("deployment.count", len(deployments)))

	var deployProtos []*deploypb.Deploy
	for _, d := range deployments {
		deployProto, err := s.deploymentToProto(ctx, repoParts[0], repoParts[1], d)
		if err != nil {
			return nil, err
		}
		deployProtos = append(deployProtos, deployProto)
	}

	return &deploypb.ListRecentDeploysResponse{
		Deploys: deployProtos,
	}, nil
}
