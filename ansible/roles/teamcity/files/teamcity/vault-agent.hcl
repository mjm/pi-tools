template {
  contents = <<EOF
{{ with secret "ssh-client-signer/sign/homelab-client" (printf "public_key=%s" (file "/opt/TeamCity/.ssh/id_rsa.pub")) "valid_principals=ubuntu,matt" }}
{{ .Data.signed_key }}
{{ end }}
EOF
  destination = "/opt/TeamCity/.ssh/signed-cert.pub"
}
