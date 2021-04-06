package deployservice

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

func (s *Server) GetDeploy(ctx context.Context, req *deploypb.GetDeployRequest) (*deploypb.GetDeployResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.Int64("deployment.id", req.GetDeployId()))

	repoParts := strings.SplitN(s.Config.GitHubRepo, "/", 2)
	d, _, err := s.GitHubClient.Repositories.GetDeployment(ctx, repoParts[0], repoParts[1], req.GetDeployId())
	if err != nil {
		return nil, err
	}

	deployResponse, err := s.deploymentToProto(ctx, repoParts[0], repoParts[1], d)
	if err != nil {
		return nil, err
	}

	return &deploypb.GetDeployResponse{
		Deploy: deployResponse,
	}, nil
}
