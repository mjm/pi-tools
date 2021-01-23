package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hashicorp/nomad/jobspec2"
)

var (
	rootPath = flag.String("root", "", "Path to resolve files from when parsing Nomad jobs")
	outPath  = flag.String("out", "", "Path to where to output resolved JSON versions of jobs")
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "error: job-resolver must be passed at least one job file")
		os.Exit(1)
	}
	jobPaths := flag.Args()
	if *rootPath == "" {
		fmt.Fprintln(os.Stderr, "error: a root path must be specified with the -root flag")
		os.Exit(1)
	}
	if *outPath == "" {
		fmt.Fprintln(os.Stderr, "error: an output path must be specified with the -out flag")
		os.Exit(1)
	}

	// first, make sure the out path exists
	if err := os.MkdirAll(*outPath, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not create output directory at %s: %s\n", *outPath, err)
		os.Exit(1)
	}

	for _, jobPath := range jobPaths {
		body, err := ioutil.ReadFile(jobPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not read job file at %s: %s\n", jobPath, err)
			os.Exit(1)
		}

		parsedJob, err := jobspec2.ParseWithConfig(&jobspec2.ParseConfig{
			Path:    jobPath,
			BaseDir: *rootPath,
			Body:    body,
			AllowFS: true,
			Strict:  true,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not parse job file at %s: %s\n", jobPath, err)
			os.Exit(1)
		}

		for _, group := range parsedJob.TaskGroups {
			for _, task := range group.Tasks {
				if task.Driver != "docker" {
					continue
				}

				configImage := task.Config["image"].(string)
				imgName, err := name.ParseReference(configImage)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: could not parse %q as a valid container image reference: %s\n", configImage, err)
					os.Exit(1)
				}

				tag, ok := imgName.(name.Tag)
				if !ok {
					continue
				}

				img, err := remote.Image(tag, remote.WithAuthFromKeychain(authn.DefaultKeychain))
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: could not resolve container image at %q: %s\n", tag, err)
					os.Exit(1)
				}

				digest, err := img.Digest()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: could not get digest of image %q: %s\n", tag, err)
					os.Exit(1)
				}

				newName := tag.Digest(digest.String())
				task.Config["image"] = newName.String()
				fmt.Fprintf(os.Stderr, "info: resolved image for task %s/%s/%s to %s\n", *parsedJob.Name, *group.Name, task.Name, newName)
			}
		}

		destPath := filepath.Join(*outPath, filepath.Base(jobPath)+".json")
		f, err := os.Create(destPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not create output file at %s: %s\n", destPath, err)
			os.Exit(1)
		}

		jsonContent, err := json.MarshalIndent(parsedJob, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: could not generate JSON for job at %s: %s\n", jobPath, err)
			os.Exit(1)
		}

		if _, err := f.Write(jsonContent); err != nil {
			fmt.Fprintf(os.Stderr, "error: could not write JSON for job at %s: %s\n", destPath, err)
			os.Exit(1)
		}

		f.Close()
		fmt.Fprintf(os.Stderr, "info: wrote job JSON to %s\n", destPath)
	}
}
