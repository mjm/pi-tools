package apiservice

import (
	"time"

	"github.com/mjm/graphql-go"
	"github.com/mjm/graphql-go/relay"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

type Deploy struct {
	*deploypb.Deploy
}

func (d *Deploy) ID() graphql.ID {
	return relay.MarshalID("deploy", d.GetId())
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
