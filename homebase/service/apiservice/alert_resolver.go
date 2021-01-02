package apiservice

import (
	"github.com/mjm/graphql-go"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type Alert struct {
	v1.Alert
}

func (a *Alert) ActiveAt() graphql.Time {
	return graphql.Time{Time: a.Alert.ActiveAt}
}
