job "vault-proxy" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "vault-proxy" {
    count = 2

    network {
      mode = "bridge"
    }

    service {
      name = "vault-proxy"
      port = 2220

      meta {
        metrics_path = "/metrics"
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
        sidecar_service {}
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
