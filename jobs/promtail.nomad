job "promtail" {
  datacenters = ["dc1"]

  type     = "system"
  priority = 80

  group "promtail" {
    network {
      mode = "bridge"
      port "expose" {}
    }

    service {
      name = "promtail"
      port = 3101

      meta {
        metrics_path = "/metrics"
        metrics_port = "${NOMAD_HOST_PORT_expose}"
      }

      connect {
        sidecar_service {
          proxy {
            expose {
              path {
                path            = "/metrics"
                protocol        = "http"
                local_path_port = 3101
                listener_port   = "expose"
              }
            }
            upstreams {
              destination_name = "loki"
              local_bind_port  = 3100
            }
          }
        }
      }
    }

    volume "run" {
      type      = "host"
      read_only = false
      source    = "promtail_run"
    }

    task "promtail" {
      driver = "docker"

      config {
        image = "grafana/promtail@sha256:d0965273b4e7c9dc2430f48e7b31f9eebf3a1d301a24c5d1cf49bdd2a9289dfb"
        args  = [
          "-config.file=${NOMAD_TASK_DIR}/promtail.yml",
        ]

        mount {
          type   = "bind"
          target = "/var/log/journal"
          source = "/var/log/journal"
        }

        mount {
          type   = "bind"
          target = "/run/log/journal"
          source = "/run/log/journal"
        }

        mount {
          type   = "bind"
          target = "/etc/machine-id"
          source = "/etc/machine-id"
        }
      }

      resources {
        cpu    = 100
        memory = 100
      }

      volume_mount {
        volume      = "run"
        destination = "/run/promtail"
        read_only   = false
      }

      template {
        data        = file("promtail/promtail.yml")
        destination = "local/promtail.yml"
      }
    }
  }
}
