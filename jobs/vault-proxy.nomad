job "vault-proxy" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "vault-proxy" {
    count = 2

    network {
      mode = "bridge"
      port "http" {
        to = 2220
      }
      port "envoy_metrics_http" {
        to = 9102
      }
    }

    service {
      name = "vault-proxy"
      port = 2220

      check {
        type                   = "http"
        port                   = "http"
        path                   = "/healthz"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
            }
          }
        }
      }
    }

    service {
      name = "deploy-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "deploy-metrics"
      port = "envoy_metrics_http"

      meta {
        metrics_path = "/metrics"
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
