package deployservice

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

var errSkipped = errors.New("skipped deploy because current version was already deployed successfully")

func (s *Server) PollForChanges(ctx context.Context, interval time.Duration) {
	s.performCheck(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.performCheck(ctx)
		}
	}
}

func (s *Server) performCheck(ctx context.Context) {
	startTime := time.Now()
	status := "succeeded"
	if err := s.checkForChanges(ctx); err != nil {
		if errors.Is(err, errSkipped) {
			status = "skipped"
		} else {
			status = "failed"
		}
	}
	duration := time.Now().Sub(startTime)

	s.deployChecksTotal.Add(ctx, 1, label.String("status", status))
	s.deployCheckDuration.Record(ctx, duration.Seconds(), label.String("status", status))
}

func (s *Server) checkForChanges(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Server.checkForChanges")
	defer span.End()

	repoParts := strings.SplitN(s.Config.GitHubRepo, "/", 2)
	deployments, _, err := s.GitHubClient.Repositories.ListDeployments(ctx, repoParts[0], repoParts[1], &github.DeploymentsListOptions{
		Task:        "deploy",
		Environment: "production",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	var prevSuccessfulCommit string

	// Assuming the first deployment returned is the most recently deployed one, so only looking at that one.
	if len(deployments) > 0 {
		prevDeploy := deployments[0]
		span.SetAttributes(
			label.Int64("deployment.previous.id", prevDeploy.GetID()),
			label.String("deployment.previous.sha", prevDeploy.GetSHA()))

		// Again, assuming the first deployment status returned is the newest one.
		statuses, _, err := s.GitHubClient.Repositories.ListDeploymentStatuses(ctx, repoParts[0], repoParts[1], prevDeploy.GetID(), &github.ListOptions{
			PerPage: 1,
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}

		if len(statuses) > 0 {
			status := statuses[0]
			span.SetAttributes(label.String("deployment.previous.status", status.GetState()))

			if status.GetState() == "success" {
				prevSuccessfulCommit = prevDeploy.GetSHA()
			}
		}
	}

	runs, _, err := s.GitHubClient.Actions.ListWorkflowRunsByFileName(ctx, repoParts[0], repoParts[1], s.Config.WorkflowName, &github.ListWorkflowRunsOptions{
		Branch: s.Config.GitHubBranch,
		Status: "success",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(label.Int("workflow.run_count", len(runs.WorkflowRuns)))
	if len(runs.WorkflowRuns) == 0 {
		return spanerr.RecordError(ctx, fmt.Errorf("no matching workflow runs found"))
	}

	run := runs.WorkflowRuns[0]
	span.SetAttributes(
		label.String("workflow.run.commit_id", run.GetHeadCommit().GetID()),
		label.String("workflow.run.commit_message", run.GetHeadCommit().GetMessage()),
		label.String("last_successful_commit_id", prevSuccessfulCommit))

	if run.GetHeadCommit().GetID() == prevSuccessfulCommit {
		span.SetAttributes(label.Bool("deploy_skipped", true))
		return errSkipped
	}

	span.SetAttributes(label.Bool("deploy_skipped", false))

	// This will get set to "success" at the end once we are done with the deploy.
	// If we don't make it that far, we'll return an error, and our defer should automatically add a "failure" status
	// to the deployment.
	finalDeploymentStatus := "failure"

	if !s.Config.DryRun {
		// First, create a new deployment for this commit
		deploy, _, err := s.GitHubClient.Repositories.CreateDeployment(ctx, repoParts[0], repoParts[1], &github.DeploymentRequest{
			Ref:              run.GetHeadCommit().ID,
			Task:             github.String("deploy"),
			AutoMerge:        github.Bool(false),
			Description:      github.String("Deploy triggered by deploy-srv"),
			RequiredContexts: &[]string{},
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}
		span.SetAttributes(
			label.Int64("deployment.id", deploy.GetID()),
			label.String("deployment.sha", deploy.GetSHA()))

		// Now set the new deployment to be in progress
		inProgressStatus, _, err := s.GitHubClient.Repositories.CreateDeploymentStatus(ctx, repoParts[0], repoParts[1], deploy.GetID(), &github.DeploymentStatusRequest{
			State:        github.String("in_progress"),
			AutoInactive: github.Bool(true),
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}
		span.SetAttributes(label.Int64("deployment.in_progress_status_id", inProgressStatus.GetID()))

		defer func(deployID int64) {
			span.SetAttributes(label.String("deployment.status", finalDeploymentStatus))

			// Create the final status for the deployment
			finalStatus, _, err := s.GitHubClient.Repositories.CreateDeploymentStatus(ctx, repoParts[0], repoParts[1], deployID, &github.DeploymentStatusRequest{
				State:        &finalDeploymentStatus,
				AutoInactive: github.Bool(true),
			})
			if err != nil {
				span.RecordError(err)
			}

			span.SetAttributes(label.Int64("deployment.final_status_id", finalStatus.GetID()))
		}(deploy.GetID())
	}

	// Alright, on with the show.
	artifacts, _, err := s.GitHubClient.Actions.ListWorkflowRunArtifacts(ctx, repoParts[0], repoParts[1], run.GetID(), nil)
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(label.Int("workflow.run.artifact_count", len(artifacts.Artifacts)))

	var artifact *github.Artifact
	for _, a := range artifacts.Artifacts {
		if a.GetName() == s.Config.ArtifactName {
			artifact = a
			break
		}
	}

	if artifact == nil {
		return spanerr.RecordError(ctx, fmt.Errorf("no artifact found named %q", s.Config.ArtifactName))
	}
	span.SetAttributes(label.Int64("workflow.run.artifact_id", artifact.GetID()))

	downloadURL, _, err := s.GitHubClient.Actions.DownloadArtifact(ctx, repoParts[0], repoParts[1], artifact.GetID(), true)
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL.String(), nil)
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	r, w := io.Pipe()
	var buf bytes.Buffer
	go buf.ReadFrom(r)

	if _, err := s.GitHubClient.Do(ctx, req, w); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(label.Int("archive.length", buf.Len()))

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(label.Int("archive.file_count", len(zipReader.File)))

	for _, f := range zipReader.File {
		if f.Name != s.Config.FileToApply+".yaml" {
			continue
		}

		span.SetAttributes(label.String("archive.file.name", f.Name))

		rc, err := f.Open()
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}

		if err := s.applyKubernetesResources(ctx, rc); err != nil {
			return spanerr.RecordError(ctx, err)
		}

		finalDeploymentStatus = "success"
		return nil
	}

	return spanerr.RecordError(ctx, fmt.Errorf("no file found in archive named %s.yaml", s.Config.FileToApply))
}

func (s *Server) applyKubernetesResources(ctx context.Context, r io.ReadCloser) error {
	ctx, span := tracer.Start(ctx, "Server.applyKubernetesResources")
	defer span.End()

	defer r.Close()

	span.SetAttributes(label.Bool("dry_run", s.Config.DryRun))
	if s.Config.DryRun {
		return nil
	}

	cmd := exec.Command("kubectl", "apply", "--server-side", "-f", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	var copyErr error
	go func() {
		defer stdin.Close()
		_, copyErr = io.Copy(stdin, r)
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s", out)
		return spanerr.RecordError(ctx, err)
	}

	if copyErr != nil {
		return spanerr.RecordError(ctx, copyErr)
	}

	return nil
}
