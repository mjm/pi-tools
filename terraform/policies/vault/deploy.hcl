# Allow reading the GitHub PAT used to check for builds and create/update deployments
path "kv/data/deploy" {
  capabilities = ["read"]
}

# Allow sending push notifications about deploys
path "kv/data/pushover" {
  capabilities = ["read"]
}

# Allow issuing client certs for accessing the Nomad API over mTLS
path "pki-int/issue/nomad-cluster" {
  capabilities = ["update"]
}

# Allow submitting jobs to Nomad
path "nomad/creds/deploy" {
  capabilities = ["read"]
}

# Allow using Consul for leader-election
path "consul/creds/deploy" {
  capabilities = ["read"]
}
