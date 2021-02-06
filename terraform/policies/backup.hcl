# Allow the Tarsnap backup job to read the Tarsnap key
path "kv/data/tarsnap" {
  capabilities = ["read"]
}

# Allow both backup jobs to get a Consul token that allows them to save a snapshot
path "consul/creds/backup" {
  capabilities = ["read"]
}
