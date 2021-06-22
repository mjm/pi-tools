package nomadic

import (
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
