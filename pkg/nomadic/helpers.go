package nomadic

import (
	nomadapi "github.com/hashicorp/nomad/api"
)

func ConsulUpstream(dest string, port int) *nomadapi.ConsulUpstream {
	return &nomadapi.ConsulUpstream{
		DestinationName: dest,
		LocalBindPort:   port,
	}
}
