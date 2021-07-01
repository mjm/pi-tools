package nomadic

import (
	nomadapi "github.com/hashicorp/nomad/api"
)

var DefaultDatacenters = []string{"dc1"}

var DefaultDNS = &nomadapi.DNSConfig{
	Servers: []string{
		"10.0.2.101",
	},
}
