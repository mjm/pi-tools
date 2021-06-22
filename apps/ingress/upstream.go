package ingress

import (
	"sort"
)

type upstream struct {
	Name        string
	ServiceName string
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
	sort.Slice(upstreams, func(i, j int) bool {
		return upstreams[i].Name < upstreams[j].Name
	})
	return upstreams
}
