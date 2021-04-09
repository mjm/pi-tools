job "promtail" {
  datacenters = ["dc1"]

  type     = "system"
  priority = 80

  group "promtail" {
    network {
      mode = "bridge"
    }

    service {
      name = "promtail"
      port = 3101

      meta {
        metrics_path = "/metrics"
      }

      check {
        type     = "http"
        expose   = true
        path     = "/ready"
        interval = "15s"
        timeout  = "3s"
      }

      connect {
        sidecar_service {
          proxy {
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
        # promtail 2.2.1
        image = "grafana/promtail@sha256:ca2711bece9b74ce51aad398dedeba706c553f16446a79d0b495573a0060529b"
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
        memory = 50
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
