# Allow the consul-exporter to get a Consul token for reading service health
path "consul/creds/prometheus" {
  capabilities = ["read"]
}
