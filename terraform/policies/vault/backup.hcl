# Allow both backup jobs to get a Consul token that allows them to save a snapshot
path "consul/creds/backup" {
  capabilities = ["read"]
}
