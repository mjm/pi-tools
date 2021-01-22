# Allow alertmanager to read Pushover secrets for sending notifications
path "kv/data/alertmanager/*" {
  capabilities = ["read"]
}
