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
        OAUTH2_PROXY_CLIENT_ID = "18da98312638f7ea59f2"
        VAULT_ADDR             = "http://active.vault.service.consul:8200"
      }

      vault {
        policies = ["oauth-proxy"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/oauth-proxy" }}
OAUTH2_PROXY_CLIENT_SECRET={{ .Data.data.client_secret }}
OAUTH2_PROXY_COOKIE_SECRET={{ .Data.data.cookie_secret }}
{{ end }}
EOF
        destination = "secrets/proxy.env"
        env         = true
      }
    }
  }
}
