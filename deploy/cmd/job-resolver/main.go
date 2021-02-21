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
	"github.com/hashicorp/nomad/api"
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
			updateNetwork(parsedJob, group)

			for _, task := range group.Tasks {
				// set common environment variables we expect to have for tracing
				if task.Env == nil {
					task.Env = map[string]string{}
				}
				task.Env["HOSTNAME"] = "${attr.unique.hostname}"
				task.Env["NOMAD_CLIENT_ID"] = "${node.unique.id}"

				// remaining modifications are only for Docker tasks
				if task.Driver != "docker" {
					continue
				}

				updateTaskLoggingConfig(parsedJob, group, task)

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

func updateNetwork(job *api.Job, group *api.TaskGroup) {
	for i, svc := range group.Services {
		if svc.Connect == nil || svc.Connect.SidecarService == nil {
			continue
		}

		if _, ok := svc.Meta["envoy_metrics_port"]; ok {
			continue
		}

		port := 9102 + i
		portLabel := fmt.Sprintf("envoy_metrics_%d", i)
		net := group.Networks[0]
		net.DynamicPorts = append(net.DynamicPorts, api.Port{
			Label:       portLabel,
			To:          port,
			HostNetwork: "default",
		})

		sidecar := svc.Connect.SidecarService
		if sidecar.Proxy == nil {
			sidecar.Proxy = &api.ConsulProxy{}
		}
		if sidecar.Proxy.Config == nil {
			sidecar.Proxy.Config = map[string]interface{}{}
		}

		sidecar.Proxy.Config["envoy_prometheus_bind_addr"] = fmt.Sprintf("0.0.0.0:%d", port)

		if svc.Meta == nil {
			svc.Meta = map[string]string{}
		}
		svc.Meta["envoy_metrics_port"] = fmt.Sprintf("${NOMAD_HOST_PORT_%s}", portLabel)
	}
}

func updateTaskLoggingConfig(job *api.Job, group *api.TaskGroup, task *api.Task) {
	// inject logging configuration so that all our tasks get logged to the systemd journal
	// with an appropriate tag
	var logging map[string]interface{}
	if loggingInt, ok := task.Config["logging"]; ok {
		loggings, ok := loggingInt.([]map[string]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "error: logging config for task %s/%s/%s is a %T, not an array of maps\n", *job.Name, *group.Name, task.Name, loggingInt)
			os.Exit(1)
		}

		if len(loggings) < 1 {
			logging = map[string]interface{}{}
			task.Config["logging"] = logging
		} else {
			logging = loggings[0]
		}
	} else {
		logging = map[string]interface{}{}
		task.Config["logging"] = logging
	}

	var loggingType string
	if loggingTypeInt, ok := logging["type"]; ok {
		loggingType, ok = loggingTypeInt.(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: logging type for task %s/%s/%s is not a string\n", *job.Name, *group.Name, task.Name)
			os.Exit(1)
		}
	} else {
		loggingType = "journald"
		logging["type"] = loggingType
	}

	var loggingCfg map[string]interface{}
	if loggingCfgInt, ok := logging["config"]; ok {
		loggingCfgs, ok := loggingCfgInt.([]map[string]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "error: logging config for task %s/%s/%s is a %T, not an array of maps\n", *job.Name, *group.Name, task.Name, loggingCfgInt)
			os.Exit(1)
		}

		if len(loggingCfgs) < 1 {
			loggingCfg = map[string]interface{}{}
			logging["config"] = []map[string]interface{}{
				loggingCfg,
			}
		} else {
			loggingCfg = loggingCfgs[0]
		}
	} else {
		loggingCfg = map[string]interface{}{}
		logging["config"] = []map[string]interface{}{
			loggingCfg,
		}
	}

	if _, ok := loggingCfg["tag"]; ok {
		return
	}

	logTag := task.Name
	if jobTag, ok := job.Meta["logging_tag"]; ok {
		logTag = jobTag
	}
	if groupTag, ok := group.Meta["logging_tag"]; ok {
		logTag = groupTag
	}
	if taskTag, ok := task.Meta["logging_tag"]; ok {
		logTag = taskTag
	}
	loggingCfg["tag"] = logTag
}
