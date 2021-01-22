# Allow the Tarsnap backup job to read the Tarsnap key
path "kv/data/tarsnap" {
  capabilities = ["read"]
}
