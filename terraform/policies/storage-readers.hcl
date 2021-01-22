# Allow postgresql to read its own root user password
path "kv/data/storage/*" {
  capabilities = ["read"]
}
