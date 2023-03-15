package ingress

import (
	"sort"

	nomadapi "github.com/hashicorp/nomad/api"

	"github.com/mjm/pi-tools/pkg/nomadic"
)

var extraUpstreams = []upstream{
	{
		Name:        "homebase",
		ServiceName: "homebase",
		ConnectPort: 3001,
	},
	{
		Name:        "homebase-api",
		ServiceName: "homebase-api",
		ConnectPort: 6460,
	},
	{
		Name:        "detect-presence",
		ServiceName: "detect-presence",
		ConnectPort: 2120,
	},
}

type upstream struct {
	Name        string
	Path        string
	ServiceName string
	ServicePort int
	ConnectPort int
	IPHash      bool
	Secure      bool
}

func sortedUpstreams() []upstream {
	// if at some point we ever reuse an upstream, modify this to unique them

	var upstreams []upstream
	for _, vhost := range virtualHosts {
		upstreams = append(upstreams, vhost.Upstream)
	}
	for _, u := range extraUpstreams {
		upstreams = append(upstreams, u)
	}
	sort.Slice(upstreams, func(i, j int) bool {
		return upstreams[i].Name < upstreams[j].Name
	})
	return upstreams
}

func connectUpstreams() []*nomadapi.ConsulUpstream {
	var upstreams []*nomadapi.ConsulUpstream
	for _, u := range sortedUpstreams() {
		if u.ConnectPort != 0 {
			upstreams = append(upstreams, nomadic.ConsulUpstream(
				u.ServiceName, u.ConnectPort))
		}
	}
	return upstreams
}
