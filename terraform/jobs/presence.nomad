job "presence" {
  datacenters = [
    "dc1",
  ]

  type = "service"

  group "detect-presence" {
    count = 2

    network {
      mode = "bridge"
      port "http" {
        to = 2120
      }
      port "grpc" {
        to = 2121
      }
    }

    service {
      name = "detect-presence"
      port = 2120

      check {
        type = "http"
        port = "http"
        path = "/healthz"
        timeout = "3s"
        interval = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgresql"
              local_bind_port = 5432
            }
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port = 14268
            }
          }
        }
      }
    }

    service {
      name = "detect-presence-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
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
        image = "mmoriarity/detect-presence-srv@__DIGEST__"
        command = "/detect-presence-srv"
        args = [
          "-db",
          "dbname=presence host=localhost sslmode=disable",
          // TODO homebase-bot
          "-mode",
          "client",
        ]

        logging {
          type = "journald"
          config {
            tag = "detect-presence-srv"
          }
        }
      }

      resources {
        cpu = 50
        memory = 50
      }

      vault {
        policies = ["presence"]
      }

      template {
        data = <<EOF
{{ with secret "database/creds/presence" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env = true
        change_mode = "restart"
      }
    }
  }
}
