job "vault-proxy" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "vault-proxy" {
    count = 2

    network {
      mode = "bridge"
      port "expose" {}
    }

    service {
      name = "vault-proxy"
      port = 2220

      meta {
        metrics_path = "/metrics"
        metrics_port = "${NOMAD_HOST_PORT_expose}"
      }

      check {
        type                   = "http"
        expose                 = true
        path                   = "/healthz"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            expose {
              path {
                path            = "/metrics"
                protocol        = "http"
                local_path_port = 2220
                listener_port   = "expose"
              }
            }
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
            }
          }
        }
      }
    }

    task "vault-proxy" {
      driver = "docker"

      config {
        image   = "mmoriarity/vault-proxy"
        command = "/vault-proxy"
      }

      env {
        VAULT_ADDR = "http://active.vault.service.consul:8200"
      }

      vault {
        policies = ["vault-proxy"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/vault-proxy" }}
COOKIE_KEY={{ .Data.data.cookie_secret }}
{{ end }}
EOF
        destination = "secrets/proxy.env"
        env         = true
      }
    }
  }
}
