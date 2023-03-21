template {
  source      = "/usr/local/etc/prometheus.yml.tpl"
  destination = "/usr/local/etc/prometheus.yml"
  command     = "service prometheus reload"
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ end }}
EOF
  destination = "/usr/local/etc/nomad.crt"
  command     = "service prometheus reload"
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.private_key }}
{{ end }}
EOF
  destination = "/usr/local/etc/nomad.key"
  command     = "service prometheus reload"
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
  destination = "/usr/local/etc/nomad.ca.crt"
  command     = "service prometheus reload"
}
