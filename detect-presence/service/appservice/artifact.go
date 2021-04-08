package appservice

import (
	"context"
	"fmt"

	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) getWorkflowArtifact(ctx context.Context, name string) (*github.Artifact, error) {
	ctx, span := tracer.Start(ctx, "Server.getWorkflowArtifact",
		trace.WithAttributes(attribute.String("github.artifact.name", name)))
	defer span.End()

	runs, _, err := s.GithubClient.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, "ios.yaml", &github.ListWorkflowRunsOptions{
		Branch: "main",
		Status: "success",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return nil, spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(attribute.Int("github.run_count", len(runs.WorkflowRuns)))
	run := runs.WorkflowRuns[0]
	span.SetAttributes(
		attribute.Int64("github.run_id", run.GetID()),
		attribute.String("github.commit_id", run.HeadCommit.GetID()),
		attribute.String("github.commit_message", run.HeadCommit.GetMessage()))

	artifacts, _, err := s.GithubClient.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, run.GetID(), nil)
	if err != nil {
		return nil, spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(attribute.Int("github.artifact_count", len(artifacts.Artifacts)))

	for _, a := range artifacts.Artifacts {
		if a.GetName() == name {
			return a, nil
		}
	}

	return nil, fmt.Errorf("no artifact found named %s", name)
}
