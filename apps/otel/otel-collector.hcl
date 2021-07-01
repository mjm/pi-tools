# Allow the OpenTelemetry collector to read the Honeycomb API key
path "kv/data/honeycomb" {
  capabilities = ["read"]
}
