template {
  contents = <<EOF
export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin
export BORG_RSH="ssh -i /usr/local/homelab/id_rsa"
export BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK=yes

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

export PHX_SERVER=true
EOF
  destination = "/usr/local/homelab/.env.sh"
  perms = "0600"
  command = "service homelab restart"
}

template {
  contents = <<EOF
{{ with secret "kv/tarsnap" }}{{ .Data.data.key | base64Decode }}{{ end }}
EOF
  destination = "/usr/local/homelab/tarsnap.key"
  perms = "0600"
}

template {
  contents = <<EOF
{{ with secret "kv/borg" }}{{ .Data.data.private_key }}{{ end }}
EOF
  destination = "/usr/local/homelab/id_rsa"
  perms = "0600"
}
