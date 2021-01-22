# Allow oauth-proxy to look up GitHub OAuth secrets
path "kv/data/oauth-proxy" {
  capabilities = ["read"]
}
