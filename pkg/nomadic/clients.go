package nomadic

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	nomadapi "github.com/hashicorp/nomad/api"
	vaultapi "github.com/hashicorp/vault/api"
)

type Clients struct {
	Nomad  *nomadapi.Client
	Consul *consulapi.Client
	Vault  *vaultapi.Client
}

func DefaultClients() (Clients, error) {
	nomadClient, err := nomadapi.NewClient(nomadapi.DefaultConfig())
	if err != nil {
		return Clients{}, err
	}

	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return Clients{}, err
	}

	vaultClient, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return Clients{}, err
	}

	return Clients{
		Nomad:  nomadClient,
		Consul: consulClient,
		Vault:  vaultClient,
	}, nil
}

func (c Clients) DeployJobs(ctx context.Context, jobs ...*nomadapi.Job) error {
	var jobsToWatch []*nomadapi.Job
	for _, j := range jobs {
		job, err := c.submitNomadJob(ctx, j)
		if err != nil {
			return err
		}

		if job != nil {
			jobsToWatch = append(jobsToWatch, job)
		}
	}

	if len(jobsToWatch) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(jobsToWatch))

	wg.Add(len(jobsToWatch))
	for _, job := range jobsToWatch {
		go func(job *nomadapi.Job) {
			defer wg.Done()

			if err := c.watchJobDeployment(ctx, job); err != nil {
				errChan <- err
			}
		}(job)
	}

	wg.Wait()

	close(errChan)
	var errs []error
	var errDescs []string
	for err := range errChan {
		errs = append(errs, err)
		errDescs = append(errDescs, err.Error())
	}

	events := Events(ctx)

	if len(errs) == 0 {
		events.Info("All jobs finished deploying successfully")
		return nil
	}

	jobWord := "jobs"
	if len(errs) == 1 {
		jobWord = "job"
	}

	events.Error("%d %s failed to deploy", len(errs), jobWord,
		withDescription(strings.Join(errDescs, "\n")))
	return fmt.Errorf("1 or more jobs failed to deploy")
}

func (c Clients) submitNomadJob(ctx context.Context, job *nomadapi.Job) (*nomadapi.Job, error) {
	planResp, _, err := c.Nomad.Jobs().Plan(job, true, nil)
	if err != nil {
		return nil, fmt.Errorf("planning nomad job %s: %w", *job.ID, err)
	}

	if planResp.Diff.Type == "None" {
		return nil, nil
	}

	resp, _, err := c.Nomad.Jobs().Register(job, nil)
	if err != nil {
		return nil, fmt.Errorf("submitting nomad job %s: %w", *job.ID, err)
	}

	job.JobModifyIndex = &resp.JobModifyIndex
	Events(ctx).Info("Submitted job %s", *job.ID)
	return job, nil
}

func (c Clients) watchJobDeployment(ctx context.Context, job *nomadapi.Job) error {
	events := Events(ctx)

	var prevDeploy *nomadapi.Deployment
	var nomadIndex uint64
	for {
		q := &nomadapi.QueryOptions{}
		if nomadIndex > 0 {
			q.WaitIndex = nomadIndex
			q.WaitTime = 30 * time.Second
		}
		d, wm, err := c.Nomad.Jobs().LatestDeployment(*job.ID, q)
		if err != nil {
			return fmt.Errorf("watching %s: %w", *job.ID, err)
		}

		if d == nil {
			return nil
		}

		if d.JobSpecModifyIndex < *job.JobModifyIndex {
			time.Sleep(5 * time.Second)
			continue
		}

		nomadIndex = wm.LastIndex

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

			events.Info("%s/%s: Placed %d, Desired %d, Healthy %d, Unhealthy %d",
				*job.ID, name, tg.PlacedAllocs, tg.DesiredTotal, tg.HealthyAllocs, tg.UnhealthyAllocs)
		}

		switch d.Status {
		case "running":
			if prevDeploy == nil || prevDeploy.StatusDescription != d.StatusDescription {
				events.Info("%s: %s", *job.ID, d.StatusDescription)
			}
			prevDeploy = d
			continue
		case "successful":
			events.Info("%s: %s", *job.ID, d.StatusDescription)
			return nil
		case "failed":
			events.Error("%s: %s", *job.ID, d.StatusDescription)
			return fmt.Errorf("%s: deployment failed: %s", *job.ID, d.StatusDescription)
		default:
			return fmt.Errorf("%s: unexpected deployment status %q", *job.ID, d.Status)
		}
	}
}
