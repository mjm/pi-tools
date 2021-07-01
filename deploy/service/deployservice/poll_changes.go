package deployservice

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/go-github/v33/github"
	"github.com/gregdel/pushover"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/attribute"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/deploy/report"
	"github.com/mjm/pi-tools/pkg/nomadic/service/nomadicservice"
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
	// used to ensure that we run our deploy to completion before exiting
	s.lock.Lock()
	defer s.lock.Unlock()

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

	s.deployChecksTotal.Add(ctx, 1, attribute.String("status", status))
	s.deployCheckDuration.Record(ctx, duration.Seconds(), attribute.String("status", status))
}

func (s *Server) checkForChanges(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "Server.checkForChanges")
	defer span.End()

	repoParts := strings.SplitN(s.Config.GitHubRepo, "/", 2)
	owner, repo := repoParts[0], repoParts[1]

	deployments, _, err := s.GitHubClient.Repositories.ListDeployments(ctx, owner, repo, &github.DeploymentsListOptions{
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
			attribute.Int64("deployment.previous.id", prevDeploy.GetID()),
			attribute.String("deployment.previous.sha", prevDeploy.GetSHA()))

		// Again, assuming the first deployment status returned is the newest one.
		statuses, _, err := s.GitHubClient.Repositories.ListDeploymentStatuses(ctx, owner, repo, prevDeploy.GetID(), &github.ListOptions{
			PerPage: 1,
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}

		if len(statuses) > 0 {
			status := statuses[0]
			span.SetAttributes(attribute.String("deployment.previous.status", status.GetState()))

			if status.GetState() == "success" {
				prevSuccessfulCommit = prevDeploy.GetSHA()
			}
		}
	}

	runs, _, err := s.GitHubClient.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, s.Config.WorkflowName, &github.ListWorkflowRunsOptions{
		Branch: s.Config.GitHubBranch,
		Status: "success",
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(attribute.Int("workflow.run_count", len(runs.WorkflowRuns)))
	if len(runs.WorkflowRuns) == 0 {
		return spanerr.RecordError(ctx, fmt.Errorf("no matching workflow runs found"))
	}

	run := runs.WorkflowRuns[0]
	span.SetAttributes(
		attribute.String("workflow.run.commit_id", run.GetHeadCommit().GetID()),
		attribute.String("workflow.run.commit_message", run.GetHeadCommit().GetMessage()),
		attribute.String("last_successful_commit_id", prevSuccessfulCommit))

	if run.GetHeadCommit().GetID() == prevSuccessfulCommit {
		span.SetAttributes(attribute.Bool("deploy_skipped", true))
		return errSkipped
	}

	span.SetAttributes(attribute.Bool("deploy_skipped", false))

	var r report.Recorder
	r.SetCommitInfo(run.GetHeadCommit().GetID(), run.GetHeadCommit().GetMessage())

	var deployID int64

	// This will get set to "success" at the end once we are done with the deploy.
	// If we don't make it that far, we'll return an error, and our defer should automatically add a "failure" status
	// to the deployment.
	finalDeploymentStatus := "failure"

	if !s.Config.DryRun {
		// First, create a new deployment for this commit
		deploy, _, err := s.GitHubClient.Repositories.CreateDeployment(ctx, owner, repo, &github.DeploymentRequest{
			Ref:              run.GetHeadCommit().ID,
			Task:             github.String("deploy"),
			AutoMerge:        github.Bool(false),
			Description:      github.String("Deploy triggered by deploy-srv"),
			RequiredContexts: &[]string{},
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}

		deployID = deploy.GetID()
		span.SetAttributes(
			attribute.Int64("deployment.id", deployID),
			attribute.String("deployment.sha", deploy.GetSHA()))
		r.SetDeployID(deployID)

		// Now set the new deployment to be in progress
		inProgressStatus, _, err := s.GitHubClient.Repositories.CreateDeploymentStatus(ctx, owner, repo, deploy.GetID(), &github.DeploymentStatusRequest{
			State:        github.String("in_progress"),
			AutoInactive: github.Bool(true),
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}
		span.SetAttributes(attribute.Int64("deployment.in_progress_status_id", inProgressStatus.GetID()))

		defer func() {
			span.SetAttributes(attribute.String("deployment.status", finalDeploymentStatus))

			// Create the final status for the deployment
			finalStatus, _, err := s.GitHubClient.Repositories.CreateDeploymentStatus(ctx, owner, repo, deployID, &github.DeploymentStatusRequest{
				State:        &finalDeploymentStatus,
				AutoInactive: github.Bool(true),
				LogURL:       github.String(fmt.Sprintf("https://homebase.home.mattmoriarity.com/deploys/%d", deployID)),
			})
			if err != nil {
				span.RecordError(err)
			}

			span.SetAttributes(attribute.Int64("deployment.final_status_id", finalStatus.GetID()))

			// Send a notification about the deploy ending
			var title string
			if finalDeploymentStatus == "failure" {
				title = "Deployment failed"
			} else {
				title = "Deployment completed"
			}
			if _, err := s.Pushover.SendMessage(&pushover.Message{
				Title:   title,
				Message: run.GetHeadCommit().GetMessage(),
			}, s.Config.PushoverRecipient); err != nil {
				span.RecordError(err)
			}
		}()
	}

	defer func(r *report.Recorder) {
		reportContent, err := r.Marshal()
		if err != nil {
			span.RecordError(err)
			return
		}

		var key string
		if deployID == 0 {
			key = ksuid.New().String()
		} else {
			key = strconv.FormatInt(deployID, 10)
		}
		span.SetAttributes(
			attribute.String("report.key", key),
			attribute.String("report.bucket", s.Config.ReportBucket),
			attribute.Int("report.size", len(reportContent)))

		_, err = s.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket: &s.Config.ReportBucket,
			Key:    &key,
			Body:   bytes.NewReader(reportContent),
		})
		if err != nil {
			span.RecordError(err)
			return
		}
	}(&r)

	// Alright, on with the show.
	artifacts, _, err := s.GitHubClient.Actions.ListWorkflowRunArtifacts(ctx, owner, repo, run.GetID(), nil)
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(attribute.Int("workflow.run.artifact_count", len(artifacts.Artifacts)))

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
	span.SetAttributes(attribute.Int64("workflow.run.artifact_id", artifact.GetID()))

	downloadURL, _, err := s.GitHubClient.Actions.DownloadArtifact(ctx, repoParts[0], repoParts[1], artifact.GetID(), true)
	if err != nil {
		r.Error("Could not get download URL for build artifact").WithError(err)
		return spanerr.RecordError(ctx, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL.String(), nil)
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}

	read, write := io.Pipe()
	var buf bytes.Buffer
	go buf.ReadFrom(read)

	if _, err := s.GitHubClient.Do(ctx, req, write); err != nil {
		r.Error("Could not download build artifact").WithError(err)
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(attribute.Int("archive.length", buf.Len()))
	r.Info("Downloaded build artifact").
		WithDescription("Artifact size: %d bytes", buf.Len())

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		r.Error("Could not unzip build artifact").WithError(err)
		return spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(attribute.Int("archive.file_count", len(zipReader.File)))

	if s.Config.DryRun {
		return nil
	}

	nomadicFile, err := zipReader.File[0].Open()
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}
	defer nomadicFile.Close()

	tmpFile, err := os.CreateTemp("", "nomadic-")
	if err != nil {
		return spanerr.RecordError(ctx, err)
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, nomadicFile); err != nil {
		return spanerr.RecordError(ctx, err)
	}
	tmpFile.Close()

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return spanerr.RecordError(ctx, err)
	}

	doneCh := make(chan struct{})
	eventCh := make(chan *deploypb.ReportEvent)
	go func() {
		for evt := range eventCh {
			r.AppendEvent(evt)
		}
		close(doneCh)
	}()

	deployErr := nomadicservice.DeployAll(ctx, tmpFile.Name(), eventCh)
	<-doneCh

	if deployErr == nil {
		r.Info("All apps finished deploying successfully")
		finalDeploymentStatus = "success"
	}

	return nil
}
