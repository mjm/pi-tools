# Allow paperless jail to read credentials for accessing paperless database
path "database/creds/paperless" {
  capabilities = ["read"]
}
