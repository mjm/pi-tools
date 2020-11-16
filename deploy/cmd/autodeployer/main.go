package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

var (
	githubRepo   = flag.String("repo", "mjm/pi-tools", "GitHub repository to pull down artifact from")
	githubBranch = flag.String("branch", "main", "Git branch of builds that should be considered for deploying")
	workflowName = flag.String("workflow", "k8s.yaml", "Filename of GitHub Actions workflow to pull artifact from")
	artifactName = flag.String("artifact", "all_k8s", "Name of artifact to download from workflow run")
	fileToApply  = flag.String("f", "k8s", "Extension-less name of the YAML file to apply")

	githubUsername  = flag.String("github-user", "mjm", "Username of GitHub user to authorize as for artifact download")
	githubTokenPath = flag.String("github-token-path", "/var/secrets/github-token", "Path to file containing GitHub PAT token")
)

func main() {
	flag.Parse()

	res, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/%s/runs?branch=%s&status=success&per_page=1", *githubRepo, *workflowName, *githubBranch))
	if err != nil {
		log.Panicf("failed to request workflow info: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Panicf("unexpected status code %d for workflow runs", res.StatusCode)
	}

	var workflowRunsResp struct {
		WorkflowRuns []struct {
			ArtifactsURL string `json:"artifacts_url"`
			HeadCommit   struct {
				ID      string
				Message string
			} `json:"head_commit"`
		} `json:"workflow_runs"`
	}
	if err := json.NewDecoder(res.Body).Decode(&workflowRunsResp); err != nil {
		log.Panicf("failed to decode workflows response: %v", err)
	}

	if len(workflowRunsResp.WorkflowRuns) == 0 {
		log.Printf("no workflow runs found, nothing to do")
		return
	}

	run := workflowRunsResp.WorkflowRuns[0]
	log.Printf("Getting artifacts for successful run:\n%s %s", run.HeadCommit.ID, run.HeadCommit.Message)

	res, err = http.Get(run.ArtifactsURL)
	if err != nil {
		log.Panicf("failed to request artifact info: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Panicf("unexpected status code %d for artifact info", res.StatusCode)
	}

	var artifactsResp struct {
		Artifacts []struct {
			Name               string
			ArchiveDownloadURL string `json:"archive_download_url"`
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&artifactsResp); err != nil {
		log.Panicf("failed to decode artifacts response: %v", err)
	}

	log.Printf("Found %d artifacts for this run", len(artifactsResp.Artifacts))

	var downloadURL string
	for _, a := range artifactsResp.Artifacts {
		if a.Name == *artifactName {
			downloadURL = a.ArchiveDownloadURL
			break
		}
	}

	if downloadURL == "" {
		log.Panicf("could not find artifact named %q in run", *artifactName)
	}

	log.Printf("Downloading archive from %s", downloadURL)

	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		log.Panicf("failed to create download request: %v", err)
	}

	tokenData, err := ioutil.ReadFile(*githubTokenPath)
	if err != nil {
		log.Panicf("failed to read github token: %v", err)
	}
	req.SetBasicAuth(*githubUsername, strings.TrimSpace(string(tokenData)))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("failed to download archive: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		log.Panicf("unexpected status code %d for archive download", res.StatusCode)
	}

	archiveData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panicf("failed to read archive data: %v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(archiveData), int64(len(archiveData)))
	if err != nil {
		log.Panicf("failed to create zip reader: %v", err)
	}

	for _, f := range zipReader.File {
		if f.Name != *fileToApply+".yaml" {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			log.Panicf("failed to open %s from zip archive: %v", f.Name, err)
		}

		cmd := exec.Command("kubectl", "apply", "--server-side", "-f", "-")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Panicf("failed to create stdin pipe for kubectl: %v", err)
		}

		go func() {
			defer stdin.Close()
			_, err = io.Copy(stdin, rc)
			if err != nil {
				log.Panicf("failed to read %s from zip archive: %v", f.Name, err)
			}
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Panicf("failed to apply resources with kubectl: %v", err)
		}

		fmt.Printf("%s", out)
		return
	}

	log.Panicf("no file named %s found in archive", *fileToApply+".yaml")
}
