# Allow reading the GitHub PAT used to check for builds and create/update deployments
path "kv/data/deploy" {
  capabilities = ["read"]
}

# Allow issuing client certs for accessing the Nomad API over mTLS
path "pki-int/issue/nomad-cluster" {
  capabilities = ["update"]
}
