package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	nomadapi "github.com/hashicorp/nomad/api"
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Panicf("need at least one job JSON file as an argument")
	}

	nomadClient, err := nomadapi.NewClient(nomadapi.DefaultConfig())
	if err != nil {
		log.Panicf("creating nomad client: %v", err)
	}

	for _, jobPath := range flag.Args() {
		f, err := os.Open(jobPath)
		if err != nil {
			log.Panicf("opening job file at %s: %s", jobPath, err)
		}

		var job nomadapi.Job
		if err := json.NewDecoder(f).Decode(&job); err != nil {
			log.Panicf("decoding job JSON for file at %s: %s", jobPath, err)
		}
		f.Close()

		if _, _, err := nomadClient.Jobs().Register(&job, nil); err != nil {
			log.Panicf("registering job %s: %s", *job.Name, err)
		}

		log.Printf("registered job %s", *job.Name)
	}
}
