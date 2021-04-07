package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

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

	var jobs []*nomadapi.Job
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

		resp, _, err := nomadClient.Jobs().Plan(&job, true, nil)
		if err != nil {
			log.Panicf("planning job %s: %v", *job.Name, err)
		}

		log.Printf("planned job %s:", *job.Name)
		printDiff(resp.Diff)

		if !*dryRun && resp.Diff.Type != "None" {
			resp, _, err := nomadClient.Jobs().Register(&job, nil)
			if err != nil {
				log.Panicf("registering job %s: %s", *job.Name, err)
			}

			log.Printf("registered job %s", *job.Name)

			job.JobModifyIndex = &resp.JobModifyIndex
			jobs = append(jobs, &job)
		}
	}

	if len(jobs) == 0 {
		return
	}

	log.Printf("Watching for deployments to complete")
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		go func(job *nomadapi.Job) {
			defer wg.Done()

			var prevDeploy *nomadapi.Deployment
			var nomadIndex uint64
			for {
				q := &nomadapi.QueryOptions{}
				if nomadIndex > 0 {
					q.WaitIndex = nomadIndex
					q.WaitTime = 30 * time.Second
				}
				d, wm, err := nomadClient.Jobs().LatestDeployment(*job.ID, q)
				if err != nil {
					return
				}

				if d == nil {
					log.Printf("Job %q doesn't have deployments, considering it complete", *job.Name)
					return
				}

				if d.JobSpecModifyIndex < *job.JobModifyIndex {
					log.Printf("Job %q hasn't created a new deployment yet, trying again", *job.Name)
					time.Sleep(5 * time.Second)
					continue
				}

				nomadIndex = wm.LastIndex
				if prevDeploy == nil || prevDeploy.StatusDescription != d.StatusDescription {
					log.Printf("%s: %s", *job.Name, d.StatusDescription)
				}

				for name, tg := range d.TaskGroups {
					// Skip output if it's the same as the last time
					if prevDeploy != nil {
						prevTG := prevDeploy.TaskGroups[name]
						if prevTG.PlacedAllocs == tg.PlacedAllocs &&
							prevTG.DesiredTotal == tg.DesiredTotal &&
							prevTG.HealthyAllocs == tg.HealthyAllocs &&
							prevTG.UnhealthyAllocs == tg.UnhealthyAllocs {
							continue
						}
					}

					log.Printf("%s/%s: placed %d, desired %d, healthy %d, unhealthy %d",
						*job.Name, name, tg.PlacedAllocs, tg.DesiredTotal, tg.HealthyAllocs, tg.UnhealthyAllocs)
				}

				switch d.Status {
				case "running":
					prevDeploy = d
					continue
				case "successful", "failed":
					return
				default:
					log.Printf("Unexpected deployment status %q, giving up", d.Status)
					return
				}
			}
		}(job)
	}

	wg.Wait()
	log.Printf("All jobs finished deploying")
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
