package appservice

import (
	"io"
	"net/http"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) InstallApp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	artifact, err := s.getWorkflowArtifact(ctx, "Presence.ipa")
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(label.Int64("github.artifact_id", artifact.GetID()))

	downloadURL, _, err := s.GithubClient.Actions.DownloadArtifact(ctx, owner, repo, artifact.GetID(), true)
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL.String(), nil)
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, res.Body); err != nil {
		_ = spanerr.RecordError(ctx, err)
	}
}
