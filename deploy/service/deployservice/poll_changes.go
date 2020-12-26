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
		label.String("last_successful_commit_id", s.lastSuccessfulCommit))

	if run.GetHeadCommit().GetID() == s.lastSuccessfulCommit {
		span.SetAttributes(label.Bool("deploy_skipped", true))
		return errSkipped
	}

	span.SetAttributes(label.Bool("deploy_skipped", false))

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

		s.lastSuccessfulCommit = run.GetHeadCommit().GetID()
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
