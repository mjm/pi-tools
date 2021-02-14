# Allow alertmanager to read Pushover secrets for sending notifications
path "kv/data/pushover" {
  capabilities = ["read"]
}

path "kv/data/pagerduty" {
  capabilities = ["read"]
}
