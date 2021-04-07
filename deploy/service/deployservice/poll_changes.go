package deployservice

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/go-github/v33/github"
	"github.com/gregdel/pushover"
	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/segmentio/ksuid"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/deploy/report"
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

	s.deployChecksTotal.Add(ctx, 1, label.String("status", status))
	s.deployCheckDuration.Record(ctx, duration.Seconds(), label.String("status", status))
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
			label.Int64("deployment.previous.id", prevDeploy.GetID()),
			label.String("deployment.previous.sha", prevDeploy.GetSHA()))

		// Again, assuming the first deployment status returned is the newest one.
		statuses, _, err := s.GitHubClient.Repositories.ListDeploymentStatuses(ctx, owner, repo, prevDeploy.GetID(), &github.ListOptions{
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
			label.Int64("deployment.id", deployID),
			label.String("deployment.sha", deploy.GetSHA()))
		r.SetDeployID(deployID)

		// Now set the new deployment to be in progress
		inProgressStatus, _, err := s.GitHubClient.Repositories.CreateDeploymentStatus(ctx, owner, repo, deploy.GetID(), &github.DeploymentStatusRequest{
			State:        github.String("in_progress"),
			AutoInactive: github.Bool(true),
		})
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}
		span.SetAttributes(label.Int64("deployment.in_progress_status_id", inProgressStatus.GetID()))

		defer func() {
			span.SetAttributes(label.String("deployment.status", finalDeploymentStatus))

			// Create the final status for the deployment
			finalStatus, _, err := s.GitHubClient.Repositories.CreateDeploymentStatus(ctx, owner, repo, deployID, &github.DeploymentStatusRequest{
				State:        &finalDeploymentStatus,
				AutoInactive: github.Bool(true),
				LogURL:       github.String(fmt.Sprintf("https://homebase.homelab/deploys/%d", deployID)),
			})
			if err != nil {
				span.RecordError(err)
			}

			span.SetAttributes(label.Int64("deployment.final_status_id", finalStatus.GetID()))

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
			label.String("report.key", key),
			label.String("report.bucket", s.Config.ReportBucket),
			label.Int("report.size", len(reportContent)))

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

	span.SetAttributes(label.Int("archive.length", buf.Len()))
	r.Info("Downloaded build artifact").
		WithDescription("Artifact size: %d bytes", buf.Len())

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		r.Error("Could not unzip build artifact").WithError(err)
		return spanerr.RecordError(ctx, err)
	}

	r.Info("Found %d Nomad jobs to submit", len(zipReader.File))

	var jobsToWatch []*nomadapi.Job
	for _, file := range zipReader.File {
		job, err := s.submitNomadJob(ctx, &r, file)
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}

		if job != nil {
			jobsToWatch = append(jobsToWatch, job)
		}
	}

	if s.Config.DryRun {
		return nil
	}

	r.Info("Watching deployment progress for all jobs")

	var wg sync.WaitGroup
	errChan := make(chan error, len(jobsToWatch))

	wg.Add(len(jobsToWatch))
	for _, job := range jobsToWatch {
		go func(job *nomadapi.Job) {
			// Don't wait for the deploy job to complete, since this server running may prevent that.
			if *job.ID == "deploy" {
				wg.Done()
			} else {
				defer wg.Done()
			}

			if err := s.watchJobDeployment(ctx, &r, job); err != nil {
				errChan <- err
			}
		}(job)
	}

	wg.Wait()

	close(errChan)
	var errs []error
	var errDescs []string
	for err := range errChan {
		errs = append(errs, err)
		errDescs = append(errDescs, err.Error())
	}

	if len(errs) == 0 {
		r.Info("All jobs finished deploying successfully")
		finalDeploymentStatus = "success"
	} else {
		jobWord := "jobs"
		if len(errs) == 1 {
			jobWord = "job"
		}

		r.Error("%d %s failed to deploy", len(errs), jobWord).
			WithDescription(strings.Join(errDescs, "\n"))
	}

	return nil
}

func (s *Server) submitNomadJob(ctx context.Context, r *report.Recorder, file *zip.File) (*nomadapi.Job, error) {
	ctx, span := tracer.Start(ctx, "Server.submitNomadJob",
		trace.WithAttributes(
			label.String("job.filename", file.Name),
			label.Uint64("job.size", file.UncompressedSize64)))
	defer span.End()

	f, err := file.Open()
	if err != nil {
		r.Error("Could not read job contents for %q", file.Name).WithError(err)
		return nil, spanerr.RecordError(ctx, err)
	}
	defer f.Close()

	// first, parse the JSON into a job
	var job nomadapi.Job
	if err := json.NewDecoder(f).Decode(&job); err != nil {
		r.Error("Could not decode job contents for %q", file.Name).WithError(err)
		return nil, spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(
		label.String("job.id", *job.ID),
		label.String("job.name", *job.Name))

	planResp, _, err := s.NomadClient.Jobs().Plan(&job, true, nil)
	if err != nil {
		r.Error("Could not plan job %q", *job.Name).WithError(err)
		return nil, spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(
		label.String("plan.diff_type", planResp.Diff.Type))

	if planResp.Diff.Type == "None" {
		return nil, nil
	}

	if s.Config.DryRun {
		r.Info("Skipped submitting job %q because this is a dry-run", *job.Name)
		return &job, nil
	}

	resp, _, err := s.NomadClient.Jobs().Register(&job, nil)
	if err != nil {
		r.Error("Could not submit job %q", *job.Name).WithError(err)
		return nil, spanerr.RecordError(ctx, err)
	}

	span.SetAttributes(label.Uint64("job.modify_index", resp.JobModifyIndex))
	job.JobModifyIndex = &resp.JobModifyIndex

	r.Info("Submitted job %q", *job.Name)
	return &job, nil
}

func (s *Server) watchJobDeployment(ctx context.Context, r *report.Recorder, job *nomadapi.Job) error {
	ctx, span := tracer.Start(ctx, "Server.watchJobDeployment",
		trace.WithAttributes(
			label.String("job.id", *job.ID)))
	defer span.End()

	var prevDeploy *nomadapi.Deployment
	var nomadIndex uint64
	for {
		q := &nomadapi.QueryOptions{}
		if nomadIndex > 0 {
			q.WaitIndex = nomadIndex
			q.WaitTime = 30 * time.Second
		}
		d, wm, err := s.NomadClient.Jobs().LatestDeployment(*job.ID, q)
		if err != nil {
			return spanerr.RecordError(ctx, err)
		}

		if d == nil {
			span.SetAttributes(label.Bool("job.has_deployments", false))
			span.AddEvent("deploy_update",
				trace.WithAttributes(
					label.String("deployment.status", "successful")))
			return nil
		}

		if d.JobSpecModifyIndex < *job.JobModifyIndex {
			span.AddEvent("wait_for_deployment",
				trace.WithAttributes(
					label.Uint64("job.modify_index", *job.JobModifyIndex),
					label.Uint64("deployment.job_modify_index", d.JobSpecModifyIndex)))
			time.Sleep(5 * time.Second)
			continue
		}

		if prevDeploy == nil {
			span.SetAttributes(label.Bool("job.has_deployments", true))
		}

		nomadIndex = wm.LastIndex
		if prevDeploy == nil || prevDeploy.StatusDescription != d.StatusDescription {
			span.AddEvent("deploy_update",
				trace.WithAttributes(
					label.String("deployment.status", d.Status),
					label.String("deployment.status_description", d.StatusDescription)))
		}

		for name, tg := range d.TaskGroups {
			// Skip output if it's the same as the last time
			if prevDeploy != nil {
				prevTG := prevDeploy.TaskGroups[name]
				if prevTG.PlacedAllocs == tg.PlacedAllocs &&
					prevTG.DesiredTotal == tg.DesiredTotal &&
					prevTG.HealthyAllocs == tg.HealthyAllocs &&
					prevTG.UnhealthyAllocs == tg.UnhealthyAllocs {
					continue
				}
			}

			span.AddEvent("task_group_update",
				trace.WithAttributes(
					label.String("task_group.name", name),
					label.Int("task_group.placed_allocs", tg.PlacedAllocs),
					label.Int("task_group.desired_total", tg.DesiredTotal),
					label.Int("task_group.healthy_allocs", tg.HealthyAllocs),
					label.Int("task_group.unhealthy_allocs", tg.UnhealthyAllocs)))

			r.Info("%s/%s: Placed %d, Desired %d, Healthy %d, Unhealthy %d",
				*job.Name, name, tg.PlacedAllocs, tg.DesiredTotal, tg.HealthyAllocs, tg.UnhealthyAllocs)
		}

		switch d.Status {
		case "running":
			if prevDeploy == nil || prevDeploy.StatusDescription != d.StatusDescription {
				r.Info("%s: %s", *job.Name, d.StatusDescription)
			}
			prevDeploy = d
			continue
		case "successful":
			r.Info("%s: %s", *job.Name, d.StatusDescription)
			return nil
		case "failed":
			r.Error("%s: %s", *job.Name, d.StatusDescription)
			err := fmt.Errorf("%s: deployment failed: %s", *job.Name, d.StatusDescription)
			return spanerr.RecordError(ctx, err)
		default:
			return spanerr.RecordError(ctx, fmt.Errorf("%s: unexpected deployment status %q", *job.Name, d.Status))
		}
	}
}
