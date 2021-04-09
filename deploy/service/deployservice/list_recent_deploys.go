package deployservice

import (
	"context"
	"strings"

	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"

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

	span.SetAttributes(attribute.Int("deployment.count", len(deployments)))

	deployProtos := make([]*deploypb.Deploy, len(deployments))
	grp, grpCtx := errgroup.WithContext(ctx)
	for i, d := range deployments {
		i, d := i, d
		grp.Go(func() error {
			deployProto, err := s.deploymentToProto(grpCtx, repoParts[0], repoParts[1], d)
			if err != nil {
				return err
			}

			deployProtos[i] = deployProto
			return nil
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, err
	}

	return &deploypb.ListRecentDeploysResponse{
		Deploys: deployProtos,
	}, nil
}
