# Allow grafana to read credentials for accessing its database
path "database/creds/grafana" {
  capabilities = ["read"]
}
