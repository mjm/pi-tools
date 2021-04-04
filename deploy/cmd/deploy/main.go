package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	nomadapi "github.com/hashicorp/nomad/api"
)

var (
	dryRun = flag.Bool("dry-run", false, "Whether to actually submit the jobs")
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

		if *dryRun {
			resp, _, err := nomadClient.Jobs().Plan(&job, true, nil)
			if err != nil {
				log.Panicf("planning job %s: %v", *job.Name, err)
			}

			log.Printf("planned job %s:", *job.Name)
			printDiff(resp.Diff)
		} else {
			if _, _, err := nomadClient.Jobs().Register(&job, nil); err != nil {
				log.Panicf("registering job %s: %s", *job.Name, err)
			}

			log.Printf("registered job %s", *job.Name)
		}
	}
}

func printDiff(diff *nomadapi.JobDiff) {
	if diff.Type == "None" {
		fmt.Println("No changes.")
		return
	}

	printFieldDiffs(diff.Fields, "")
	printObjectDiffs(diff.Objects, "")

	for _, tg := range diff.TaskGroups {
		if tg.Type == "None" {
			continue
		}

		fmt.Printf("%s task group %q:\n", tg.Type, tg.Name)
		printFieldDiffs(tg.Fields, "  ")
		printObjectDiffs(tg.Objects, "  ")

		for _, task := range tg.Tasks {
			if task.Type == "None" {
				continue
			}

			fmt.Printf("  %s task %q:\n", task.Type, task.Name)
			printFieldDiffs(task.Fields, "    ")
			printObjectDiffs(task.Objects, "    ")
		}
	}
}

func printObjectDiffs(diffs []*nomadapi.ObjectDiff, indent string) {
	for _, diff := range diffs {
		fmt.Printf("%s%s %s:\n", indent, diff.Type, diff.Name)
		printFieldDiffs(diff.Fields, indent+"  ")
		printObjectDiffs(diff.Objects, indent+"  ")
	}
}

func printFieldDiffs(diffs []*nomadapi.FieldDiff, indent string) {
	for _, diff := range diffs {
		if diff.Type == "None" {
			continue
		}

		fmt.Printf("%s%s %s: %s -> %s\n", indent, diff.Type, diff.Name, diff.Old, diff.New)
	}
}
