job "presence" {
  datacenters = ["dc1"]

  type = "service"

  group "detect-presence" {
    count = 2

    network {
      mode = "bridge"
    }

    service {
      name = "detect-presence"
      port = 2120

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
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "homebase-bot-grpc"
              local_bind_port  = 6361
            }
          }
        }
      }
    }

    service {
      name = "detect-presence-grpc"
      port = 2121

      connect {
        sidecar_service {}
      }
    }

    task "detect-presence-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/detect-presence-srv"
        command = "/detect-presence-srv"
        args    = [
          "-db",
          "dbname=presence host=10.0.2.102 sslmode=disable",
          "-mode",
          "client",
        ]
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["presence"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "kv/deploy" }}{{ .Data.data.github_token }}{{ end }}
EOF
        destination = "secrets/github-token"
        change_mode = "restart"
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/presence" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
        change_mode = "restart"
      }
    }
  }
}
