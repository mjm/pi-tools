job "promtail" {
  datacenters = [
    "dc1",
  ]

  type = "system"
  priority = 80

  group "promtail" {
    network {
      mode = "bridge"
      port "http" {
        to = 3101
      }
    }

    service {
      name = "promtail"
      port = 3101

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "loki"
              local_bind_port = 3100
            }
          }
        }
      }
    }

    service {
      name = "promtail-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    volume "run" {
      type = "host"
      read_only = false
      source = "promtail_run"
    }

    task "promtail" {
      driver = "docker"

      config {
        image = "grafana/promtail@sha256:d0965273b4e7c9dc2430f48e7b31f9eebf3a1d301a24c5d1cf49bdd2a9289dfb"
        args = [
          "-config.file=${NOMAD_TASK_DIR}/promtail.yml",
        ]

        logging {
          type = "journald"
          config {
            tag = "promtail"
          }
        }

        mount {
          type = "bind"
          target = "/var/log/journal"
          source = "/var/log/journal"
        }

        mount {
          type = "bind"
          target = "/run/log/journal"
          source = "/run/log/journal"
        }

        mount {
          type = "bind"
          target = "/etc/machine-id"
          source = "/etc/machine-id"
        }
      }

      resources {
        cpu = 100
        memory = 100
      }

      volume_mount {
        volume = "run"
        destination = "/run/promtail"
        read_only = false
      }

      template {
        data = <<EOF
server:
  http_listen_port: 3101

positions:
  filename: /run/promtail/positions.yaml

clients:
  - url: http://localhost:3100/loki/api/v1/push

scrape_configs:
  - job_name: systemd-journal
    journal:
      labels:
        job: default/systemd-journal
      path: /var/log/journal
    relabel_configs:
      - source_labels:
          - __journal__systemd_unit
        target_label: systemd_unit
      - source_labels:
          - __journal__hostname
        target_label: hostname
      - source_labels:
          - __journal_syslog_identifier
        target_label: syslog_identifier
EOF
        destination = "local/promtail.yml"
      }
    }
  }
}
