# Allow vault-proxy to look up its cookie secret for sessions
path "kv/data/vault-proxy" {
  capabilities = ["read"]
}
