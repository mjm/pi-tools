# Allow detect-presence-srv to read credentials for accessing its database
path "database/creds/presence" {
  capabilities = ["read"]
}

# Allow reading the GitHub PAT used to check for builds of the iOS app
path "kv/data/deploy" {
  capabilities = ["read"]
}
