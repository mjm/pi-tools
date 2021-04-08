package appservice

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
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

	span.SetAttributes(attribute.Int64("github.artifact_id", artifact.GetID()))

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

	var buf bytes.Buffer
	if _, err := s.GithubClient.Do(ctx, req, &buf); err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("archive.length", buf.Len()))

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zipFile := zipReader.File[0]
	f, err := zipFile.Open()
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, f); err != nil {
		_ = spanerr.RecordError(ctx, err)
	}
}
