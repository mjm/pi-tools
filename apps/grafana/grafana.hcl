# Allow grafana to read credentials for accessing its database
path "database/creds/grafana" {
  capabilities = ["read"]
}

# Allow grafana to read the auth token for Fly.io to connect to Prometheus
path "kv/data/grafana" {
  capabilities = ["read"]
}
