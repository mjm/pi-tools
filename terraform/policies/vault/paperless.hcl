# Allow paperless jail to read credentials for accessing paperless database
path "database/creds/paperless" {
  capabilities = ["read"]
}

# Allow paperless jail to read the secret key
path "kv/data/paperless" {
  capabilities = ["read"]
}
