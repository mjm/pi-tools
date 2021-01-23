locals {
  dashboards = fileset(".", "grafana/dashboards/*.json")
}

job "grafana" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "grafana" {
    count = 3

    network {
      mode = "bridge"
      port "http" {
        to = 3000
      }
      port "envoy_metrics" {
        to = 9102
      }
    }

    service {
      name = "grafana"
      port = 3000

      check {
        expose   = true
        path     = "/api/health"
        type     = "http"
        interval = "15s"
        timeout  = "3s"
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgresql"
              local_bind_port  = 5432
            }
            upstreams {
              destination_name = "loki"
              local_bind_port  = 3100
            }
          }
        }
      }
    }

    service {
      name = "grafana-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "grafana-metrics"
      port = "envoy_metrics"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "grafana" {
      driver = "docker"

      config {
        image = "grafana/grafana@sha256:f0817ecbf8dcf33e10cca2245bd25439433c441189bbe1ce935ac61d05f9cc6f"

        logging {
          type = "journald"
          config {
            tag = "grafana"
          }
        }
      }

      resources {
        cpu    = 100
        memory = 150
      }

      env {
        GF_PATHS_CONFIG       = "${NOMAD_SECRETS_DIR}/grafana.ini"
        GF_PATHS_PROVISIONING = "${NOMAD_TASK_DIR}/provisioning"
      }

      vault {
        policies = ["grafana"]
      }

      template {
        data        = file("grafana/grafana.ini")
        destination = "secrets/grafana.ini"
        change_mode = "restart"
      }

      template {
        data        = file("grafana/datasources.yaml")
        destination = "local/provisioning/datasources/datasources.yaml"
        change_mode = "restart"
      }

      template {
        data        = file("grafana/dashboards.yaml")
        destination = "local/provisioning/dashboards/dashboards.yaml"
      }

      dynamic "template" {
        for_each = local.dashboards

        content {
          data = file(template.value)
          destination = "local/dashboards/${basename(template.value)}"
          // prevent interpreting blocks delimited by '{{' and '}}' as consul templates
          left_delimiter = "do_not_substitute"
        }
      }
    }
  }
}
