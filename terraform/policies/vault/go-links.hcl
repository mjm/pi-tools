# Allow go-links to read credentials for accessing its database
path "database/creds/go-links" {
  capabilities = ["read"]
}
