job "vault-proxy" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "vault-proxy" {
    count = 2

    network {
      port "http" {
        to = "2220"
      }
    }

    service {
      name = "vault-proxy"
      port = "http"

      check {
        type     = "http"
        path     = "/healthz"
        interval = "15s"
        timeout  = "3s"
      }
    }

    task "vault-proxy" {
      driver = "docker"

      config {
        image   = "mmoriarity/vault-proxy"
        command = "/vault-proxy"
        ports   = ["http"]
      }

      env {
        VAULT_ADDR = "http://active.vault.service.consul:8200"
      }

      vault {
        policies = ["oauth-proxy"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/oauth-proxy" }}
COOKIE_KEY={{ .Data.data.cookie_secret }}
{{ end }}
EOF
        destination = "secrets/proxy.env"
        env         = true
      }
    }
  }
}
