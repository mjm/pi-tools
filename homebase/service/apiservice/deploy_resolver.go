package apiservice

import (
	"context"
	"strconv"
	"time"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

type Deploy struct {
	*deploypb.Deploy
	r *Resolver
}

func (d *Deploy) ID() graphql.ID {
	return relay.MarshalID("deploy", d.GetId())
}

func (d *Deploy) RawID() graphql.ID {
	return graphql.ID(strconv.FormatInt(d.GetId(), 10))
}

func (d *Deploy) CommitSHA() string {
	return d.GetCommitSha()
}

func (d *Deploy) CommitMessage() string {
	return d.GetCommitMessage()
}

func (d *Deploy) State() string {
	return d.GetState().String()
}

func (d *Deploy) StartedAt() (graphql.Time, error) {
	t, err := time.Parse(time.RFC3339, d.GetStartedAt())
	if err != nil {
		return graphql.Time{}, err
	}
	return graphql.Time{Time: t}, nil
}

func (d *Deploy) FinishedAt() (*graphql.Time, error) {
	if d.GetFinishedAt() == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, d.GetFinishedAt())
	if err != nil {
		return nil, err
	}
	return &graphql.Time{Time: t}, nil
}

func (d *Deploy) Report(ctx context.Context) (*DeployReport, error) {
	res, err := d.r.deployClient.GetDeployReport(ctx, &deploypb.GetDeployReportRequest{
		DeployId: d.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &DeployReport{Report: res.GetReport()}, nil
}

type DeployReport struct {
	*deploypb.Report
}

func (r *DeployReport) ID() graphql.ID {
	return relay.MarshalID("deployreport", r.GetDeployId())
}

func (r *DeployReport) Events() []*DeployReportEvent {
	var evts []*DeployReportEvent
	for _, e := range r.GetEvents() {
		evts = append(evts, &DeployReportEvent{ReportEvent: e})
	}
	return evts
}

type DeployReportEvent struct {
	*deploypb.ReportEvent
}

func (e *DeployReportEvent) Timestamp() graphql.Time {
	return graphql.Time{Time: e.GetTimestamp().AsTime()}
}

func (e *DeployReportEvent) Level() string {
	switch e.GetLevel() {
	case deploypb.ReportEvent_INFO:
		return "INFO"
	case deploypb.ReportEvent_WARNING:
		return "WARNING"
	case deploypb.ReportEvent_ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (e *DeployReportEvent) Summary() string {
	return e.GetSummary()
}

func (e *DeployReportEvent) Description() string {
	return e.GetDescription()
}
