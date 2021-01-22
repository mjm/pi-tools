package appservice

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"howett.net/plist"

	"github.com/mjm/pi-tools/pkg/itms"
	"github.com/mjm/pi-tools/pkg/spanerr"
)

func (s *Server) InstallManifest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	artifact, err := s.getWorkflowArtifact(ctx, "Info.plist")
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(label.Int64("github.artifact_id", artifact.GetID()))
	artifactURL, _, err := s.GithubClient.Actions.DownloadArtifact(ctx, owner, repo, artifact.GetID(), true)
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, artifactURL.String(), nil)
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader, writer := io.Pipe()
	var buf bytes.Buffer
	go buf.ReadFrom(reader)

	if _, err := s.GithubClient.Do(ctx, req, writer); err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(label.Int("archive.length", buf.Len()))
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// should be only a single file here
	file, err := zipReader.File[0].Open()
	defer file.Close()

	plistBytes, err := ioutil.ReadAll(file)
	if err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var infoPlist map[string]interface{}
	if _, err := plist.Unmarshal(plistBytes, &infoPlist); err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	manifest := &itms.Manifest{
		Items: []*itms.Item{
			{
				Metadata: itms.Metadata{
					BundleIdentifier: infoPlist["CFBundleIdentifier"].(string),
					BundleVersion:    infoPlist["CFBundleVersion"].(string),
					Kind:             "software",
					Title:            infoPlist["CFBundleName"].(string),
				},
				Assets: []*itms.Asset{
					{
						Kind: "software-package",
						URL:  fmt.Sprintf("https://%s/app/install", r.Host),
					},
				},
			},
		},
	}

	if err := plist.NewEncoderForFormat(w, plist.XMLFormat).Encode(manifest); err != nil {
		err = spanerr.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
