# Allow Prometheus to scrape Vault's metrics endpoint
path "sys/metrics" {
  capabilities = ["read", "list"]
}

# Allow Prometheus to issue itself client certificates for accessing Nomad's metrics
path "pki-int/issue/nomad-cluster" {
  capabilities = ["update"]
}

# Allow Prometheus to get a Consul token for reading service configs for discovery
path "consul/creds/prometheus" {
  capabilities = ["read"]
}
