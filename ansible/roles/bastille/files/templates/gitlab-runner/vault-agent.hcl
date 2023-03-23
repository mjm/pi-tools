template {
  contents = <<EOF
{{ with secret "ssh-client-signer/sign/homelab-client" (printf "public_key=%s" (file "/usr/local/etc/gitlab-runner/id_rsa.pub")) "valid_principals=ubuntu,matt" }}
{{ .Data.signed_key }}
{{ end }}
EOF
  destination = "/usr/local/etc/gitlab-runner/signed-cert.pub"
}

template {
  contents = <<EOF
{{ with secret "kv/deploy" }}{{ .Data.data.ansible_vault_password }}{{ end }}
EOF
  destination = "/usr/local/etc/gitlab-runner/vault-password"
}
