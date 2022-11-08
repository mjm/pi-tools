path "ssh-client-signer/sign/homelab-client" {
  capabilities = ["update"]
}

# Allow reading the Ansible vault password
path "kv/data/deploy" {
  capabilities = ["read"]
}
