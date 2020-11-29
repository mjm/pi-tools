package appservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) DownloadApp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/mjm/pi-tools/actions/workflows/ios.yaml/runs?branch=main&status=success&per_page=1", nil)
	if err != nil {
		err = spanerr.RecordError(ctx, fmt.Errorf("failed to request workflow info: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		err = spanerr.RecordError(ctx, fmt.Errorf("failed to request workflow info: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		err = spanerr.RecordError(ctx, fmt.Errorf("unexpected status code %d for workflow runs", res.StatusCode))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var workflowRunsResp struct {
		WorkflowRuns []struct {
			ArtifactsURL string `json:"artifacts_url"`
			CheckSuiteID int64  `json:"check_suite_id"`
			HeadCommit   struct {
				ID      string
				Message string
			} `json:"head_commit"`
		} `json:"workflow_runs"`
	}
	if err := json.NewDecoder(res.Body).Decode(&workflowRunsResp); err != nil {
		err = spanerr.RecordError(ctx, fmt.Errorf("failed to decode workflows response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(label.Int("github.run_count", len(workflowRunsResp.WorkflowRuns)))
	if len(workflowRunsResp.WorkflowRuns) == 0 {
		http.Error(w, "no workflow runs found", http.StatusNotFound)
		return
	}

	run := workflowRunsResp.WorkflowRuns[0]
	span.SetAttributes(
		label.String("github.commit_id", run.HeadCommit.ID),
		label.String("github.commit_message", run.HeadCommit.Message),
		label.Int64("github.check_suite_id", run.CheckSuiteID))

	req, err = http.NewRequestWithContext(ctx, http.MethodGet, run.ArtifactsURL, nil)
	if err != nil {
		err = spanerr.RecordError(ctx, fmt.Errorf("failed to request artifact info: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err = s.HTTPClient.Do(req)
	if err != nil {
		err = spanerr.RecordError(ctx, fmt.Errorf("failed to request artifact info: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		err = spanerr.RecordError(ctx, fmt.Errorf("unexpected status code %d for artifact info", res.StatusCode))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var artifactsResp struct {
		Artifacts []struct {
			Name string
			ID   int64
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&artifactsResp); err != nil {
		err = spanerr.RecordError(ctx, fmt.Errorf("failed to decode artifacts response: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(label.Int("github.artifact_count", len(artifactsResp.Artifacts)))

	var artifactID int64
	for _, artifact := range artifactsResp.Artifacts {
		if artifact.Name == "Presence.ipa" {
			artifactID = artifact.ID
			break
		}
	}

	if artifactID == 0 {
		http.Error(w, "no artifact named Presence.ipa found", http.StatusNotFound)
		return
	}

	span.SetAttributes(label.Int64("github.artifact_id", artifactID))

	downloadURL := fmt.Sprintf("https://github.com/mjm/pi-tools/suites/%d/artifacts/%d", run.CheckSuiteID, artifactID)
	http.Redirect(w, r, downloadURL, http.StatusTemporaryRedirect)
}
