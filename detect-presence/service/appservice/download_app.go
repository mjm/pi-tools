package appservice

import (
	"net/http"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"

	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) DownloadApp(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, downloadURL.String(), http.StatusTemporaryRedirect)
}
