# Allow adminer to read credentials for any database
path "database/creds/*" {
  capabilities = ["read"]
}
