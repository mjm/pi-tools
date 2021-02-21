locals {
  oauth_locations_snippet = file("ingress/oauth-locations.conf")
  oauth_request_snippet   = file("ingress/oauth-request.conf")

  // .homelab certificates that need to be issued from Vault
  homelab_certs = [
    "alertmanager",
    "auth",
    "consul",
    "go",
    "grafana",
    "homebase",
    "jaeger",
    "nomad",
    "prometheus",
    "pihole",
    "vault",
  ]
}

job "ingress" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  update {
    max_parallel = 1
    stagger      = "30s"
  }

  group "ingress" {
    count = 2

    network {
      mode = "bridge"
      port "http" {
        to     = 80
        static = 80
      }
      port "https" {
        to     = 443
        static = 443
      }
      port "envoy_metrics" {
        to = 9102
      }
    }

    service {
      name = "ingress-http"
      port = 80

      meta {
        envoy_metrics_port = "${NOMAD_HOST_PORT_envoy_metrics}"
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "detect-presence"
              local_bind_port  = 2120
            }
            upstreams {
              destination_name = "vault-proxy"
              local_bind_port  = 2220
            }
            upstreams {
              destination_name = "go-links"
              local_bind_port  = 4240
            }
            upstreams {
              destination_name = "homebase-api"
              local_bind_port  = 6460
            }
            upstreams {
              destination_name = "grafana"
              local_bind_port  = 3000
            }
          }
        }
      }
    }

    service {
      name = "ingress-https"
      port = 443
    }

    task "nginx" {
      driver = "docker"

      config {
        image   = "nginx@sha256:763d95e3db66d9bd1bb926c029e5659ee67eb49ff57f83d331de5f5af6d2ae0c"
        volumes = [
          "local:/etc/nginx/conf.d",
          "secrets:/etc/nginx/ssl",
        ]
      }

      meta {
        logging_tag = "ingress"
      }

      vault {
        policies    = ["ingress"]
        change_mode = "noop"
      }

      template {
        // unfortunately, nomad doesn't seem to support templatefile yet, so we resort to this for now
        data          = replace(replace(file("ingress/load-balancer.conf"), "__OAUTH_LOCATIONS_SNIPPET__", local.oauth_locations_snippet), "__OAUTH_REQUEST_SNIPPET__", local.oauth_request_snippet)
        destination   = "local/load-balancer.conf"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      dynamic "template" {
        for_each = local.homelab_certs

        content {
          data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=${template.value}.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
          destination   = "secrets/${template.value}.homelab.pem"
          change_mode   = "signal"
          change_signal = "SIGHUP"
        }
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/nomad.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
        destination   = "secrets/nomad.ca.crt"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }
  }
}
