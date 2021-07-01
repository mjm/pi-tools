# Allow Loki to read its password for storing logs in Minio
path "kv/data/loki" {
  capabilities = ["read"]
}
