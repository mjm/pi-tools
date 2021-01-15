job "go-links" {
  datacenters = [
    "dc1"
  ]

  type = "service"

  group "go-links" {
    count = 2

    network {
      mode = "bridge"
      port "expose_http" {}
    }

    service {
      name = "go-links"
      port = 4240

      check {
        expose = true
        type = "http"
        path = "/healthz"
        timeout = "3s"
        interval = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            expose {
              path {
                path = "/metrics"
                protocol = "http"
                local_path_port = 4240
                listener_port = "expose_http"
              }
            }
            upstreams {
              destination_name = "postgresql"
              local_bind_port = 5432
            }
          }
        }
      }
    }

    service {
      name = "go-links-grpc"
      port = 4241

//      check {
//        type = "grpc"
//        timeout = "3s"
//        interval = "15s"
//        success_before_passing = 3
//      }
      connect {
        sidecar_service {}
      }
    }

    task "go-links-srv" {
      driver = "docker"

      config {
        image = "mmoriarity/go-links-srv@__DIGEST__"
        command = "/go-links"
        args = [
          "-db",
          "dbname=golinks host=localhost sslmode=disable",
        ]
      }

      resources {
        cpu = 50
        memory = 50
      }

      vault {
        policies = ["go-links"]
      }

      template {
        data = <<EOF
{{ with secret "database/creds/go-links" }}
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
