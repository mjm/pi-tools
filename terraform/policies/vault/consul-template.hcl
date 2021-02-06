# consul-template policy allows access to whatever secrets are needed for consul-template
# running on the homelab cluster nodes.

# allow getting certificates for nomad
path "pki-int/issue/nomad-cluster" {
  capabilities = ["update"]
}

# allow getting a Consul token for nomad
path "consul/creds/nomad-client-server" {
  capabilities = ["read"]
}

# allow renewing the token we assigned, so that as long as the node doesn't stop running for several hours
# we shouldn't need to replace the token
path "auth/token/renew-self" {
  capabilities = ["update"]
}
