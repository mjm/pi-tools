# Allow homebase-bot to read credentials for accessing its database
path "database/creds/homebase-bot" {
  capabilities = ["read"]
}

# Allow homebase-bot to read its Telegram API token
path "kv/data/homebase-bot" {
  capabilities = ["read"]
}
