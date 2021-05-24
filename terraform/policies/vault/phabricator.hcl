# Allow phabricator jail and backup jobs to read credentials for accessing Phabricator databases
path "database/creds/phabricator" {
  capabilities = ["read"]
}
