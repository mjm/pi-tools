job "go-links" {
  datacenters = ["dc1"]

  type = "service"

  group "go-links" {
    count = 2

    network {
      mode = "bridge"
      port "expose" {}
      port "envoy_metrics_http" {
        to = 9102
      }
      port "envoy_metrics_grpc" {
        to = 9103
      }
    }

    service {
      name = "go-links"
      port = 4240

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
                local_path_port = 4240
                listener_port   = "expose"
              }
            }
            upstreams {
              destination_name = "postgresql"
              local_bind_port  = 5432
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
      name = "go-links-metrics"
      port = "envoy_metrics_http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "go-links-metrics"
      port = "envoy_metrics_grpc"

      meta {
        metrics_path = "/metrics"
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
        sidecar_service {
          proxy {
            config {
              envoy_prometheus_bind_addr = "0.0.0.0:9103"
            }
          }
        }
      }
    }

    task "go-links-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/go-links-srv"
        command = "/go-links"
        args    = [
          "-db",
          "dbname=golinks host=127.0.0.1 sslmode=disable",
        ]
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["go-links"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/go-links" }}
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
