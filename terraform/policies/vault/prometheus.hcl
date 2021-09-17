# Allow Prometheus to scrape Vault's metrics endpoint
path "sys/metrics" {
  capabilities = ["read", "list"]
}

# Allow Prometheus to issue itself client certificates for accessing Nomad's metrics
path "pki-int/issue/nomad-cluster" {
  capabilities = ["update"]
}
