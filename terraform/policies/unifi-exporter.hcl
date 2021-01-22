# Allow unifi-exporter to read the password to the Unifi account
path "kv/data/unifi-exporter" {
  capabilities = ["read"]
}
