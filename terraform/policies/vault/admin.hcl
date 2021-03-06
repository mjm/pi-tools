# Read system health check
path "sys/health" {
  capabilities = ["read", "sudo"]
}

# Create and manage ACL policies broadly across Vault

# List existing policies
path "sys/policies/acl" {
  capabilities = ["list"]
}

# Create and manage ACL policies
path "sys/policies/acl/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Enable and manage authentication methods broadly across Vault

# Manage auth methods broadly across Vault
path "auth/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Create, update, and delete auth methods
path "sys/auth/*" {
  capabilities = ["create", "update", "delete", "sudo"]
}

# List auth methods
path "sys/auth" {
  capabilities = ["read"]
}

# List, create, update, and delete key/value secrets
path "kv/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "database/roles/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "database/config/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "database/roles/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "ssh-client-signer/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "ssh-host-signer/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "pki/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "pki-int/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "pki-homelab/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "identity/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "consul/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

path "nomad/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# Manage secrets engines
path "sys/mounts/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}

# List existing secrets engines.
path "sys/mounts" {
  capabilities = ["read"]
}

# Manage plugin catalog
path "sys/plugins/catalog/*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
