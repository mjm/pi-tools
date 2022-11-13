template {
  contents = <<EOF
{{ with secret "auth/token/lookup-self" }}
vault {
  token = {{ .Data.id | toJSON }}
}
{{ end }}
EOF
  destination = "/usr/local/etc/consul-template.d/vault.hcl"
  command = "pkill -HUP consul-template"
}
