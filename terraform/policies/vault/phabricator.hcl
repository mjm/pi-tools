# Allow phabricator jail and backup jobs to read credentials for accessing Phabricator databases
path "database/creds/phabricator" {
  capabilities = ["read"]
}

# Allow reading fastmail password and minio secret key for Phabricator config
path "kv/data/phabricator" {
  capabilities = ["read"]
}
