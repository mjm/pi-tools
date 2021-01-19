vault {
  address     = "http://active.vault.service.consul:8200"
  renew_token = true
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "common_name=server.global.nomad" "ttl=24h" "alt_names=localhost,nomad.service.consul" "ip_sans=127.0.0.1"}}
{{ .Data.certificate }}
{{ end }}
EOF
  destination = "/etc/nomad/agent.crt"
  perms       = 0600
  command     = "systemctl reload nomad"
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "common_name=server.global.nomad" "ttl=24h" "alt_names=localhost,nomad.service.consul" "ip_sans=127.0.0.1"}}
{{ .Data.private_key }}
{{ end }}
EOF
  destination = "/etc/nomad/agent.key"
  perms       = 0600
  command     = "systemctl reload nomad"
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "common_name=server.global.nomad" "ttl=24h"}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
  destination = "/etc/nomad/ca.crt"
  command     = "systemctl reload nomad"
}


template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h"}}
{{ .Data.certificate }}
{{ end }}
EOF
  destination = "/etc/nomad/cli.crt"
}

template {
  contents    = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h"}}
{{ .Data.private_key }}
{{ end }}
EOF
  destination = "/etc/nomad/cli.key"
}
