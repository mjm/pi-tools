template {
  contents = <<EOF
{{ with secret "database/creds/homelab" -}}
export TRIPS_DATABASE_URL=postgresql://{{ .Data.username }}:{{ .Data.password }}@postgresql.service.consul/trips
export GO_LINKS_DATABASE_URL=postgresql://{{ .Data.username }}:{{ .Data.password }}@postgresql.service.consul/go_links
{{ end -}}
export PAPERLESS_TOKEN={{ with secret "kv/paperless/client" }}{{ .Data.data.api_token }}{{ end }}
export TELEGRAM_TOKEN={{ with secret "kv/homebase-bot" }}{{ .Data.data.telegram_token }}{{ end }}
{{ with secret "kv/homelab" -}}
export GITHUB_TOKEN={{ .Data.data.github_token }}
export SECRET_KEY_BASE={{ .Data.data.secret_key_base }}
{{ end }}
EOF
  destination = "/usr/local/homelab/.env.sh"
  perms = "0600"
  command = "service homelab restart"
}
