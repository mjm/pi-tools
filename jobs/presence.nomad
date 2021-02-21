job "presence" {
  datacenters = ["dc1"]

  type = "service"

  group "detect-presence" {
    count = 2

    network {
      mode = "bridge"
      port "expose" {}
    }

    service {
      name = "detect-presence"
      port = 2120

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
                local_path_port = 2120
                listener_port   = "expose"
              }
            }
            upstreams {
              destination_name = "postgresql"
              local_bind_port  = 5432
            }
            upstreams {
              destination_name = "homebase-bot-grpc"
              local_bind_port  = 6361
            }
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
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
          "dbname=presence host=127.0.0.1 sslmode=disable",
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
