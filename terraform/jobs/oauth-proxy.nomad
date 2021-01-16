job "oauth-proxy" {
  datacenters = [
    "dc1",
  ]

  type = "service"
  priority = 70

  group "oauth-proxy" {
    count = 3

    network {
      port "http" {
        to = "4180"
      }
    }

    service {
      name = "oauth-proxy"
      port = "http"

      check {
        type = "http"
        path = "/ping"
        interval = "15s"
        timeout = "3s"
      }
    }

    task "oauth-proxy" {
      driver = "docker"

      config {
        image = "quay.io/oauth2-proxy/oauth2-proxy@sha256:cf9c36686ae737ffcfe0f91e8ec60988695a0fa83748a96e079c3a0cf0a985fc"
        args = [
          "--provider=github",
          "--email-domain=*",
          "--github-org=mmoriarity",
          "--upstream=file:///dev/null",
          "--http-address=0.0.0.0:4180",
          "--cookie-secure=false",
          "--cookie-domain=.homelab",
          "--redirect-url=https://homebase.homelab/oauth2/callback",
          "--set-xauthrequest",
          "--pass-access-token",
        ]
        ports = ["http"]

        logging {
          type = "journald"
          config {
            tag = "oauth-proxy"
          }
        }
      }

      env {
        OAUTH2_PROXY_CLIENT_ID = "18da98312638f7ea59f2"
      }

      vault {
        policies = ["oauth-proxy"]
      }

      template {
        data = <<EOF
{{ with secret "kv/oauth-proxy" }}
OAUTH2_PROXY_CLIENT_SECRET={{ .Data.data.client_secret }}
OAUTH2_PROXY_COOKIE_SECRET={{ .Data.data.cookie_secret }}
{{ end }}
EOF
        destination = "secrets/proxy.env"
        env = true
      }
    }
  }
}
