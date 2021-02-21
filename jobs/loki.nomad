job "loki" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 80

  group "loki" {
    count = 1

    network {
      mode = "bridge"
      port "expose" {}
      port "envoy_metrics" {
        to = 9102
      }
    }

    service {
      name = "loki"
      port = 3100

      meta {
        metrics_path       = "/metrics"
        metrics_port       = "${NOMAD_HOST_PORT_expose}"
        envoy_metrics_port = "${NOMAD_HOST_PORT_envoy_metrics}"
      }

      connect {
        sidecar_service {
          proxy {
            expose {
              path {
                path            = "/metrics"
                protocol        = "http"
                local_path_port = 3100
                listener_port   = "expose"
              }
            }
          }
        }
      }
    }

    task "loki" {
      driver = "docker"

      config {
        image = "grafana/loki@sha256:6afc0da6995fecf15307762d378242b65cab20d4a35b4a39397d67cad48fb7fb"
        args  = ["-config.file=${NOMAD_TASK_DIR}/loki.yml"]
      }

      resources {
        cpu    = 50
        memory = 200
      }

      template {
        data        = file("loki/loki.yml")
        destination = "local/loki.yml"
      }
    }
  }
}
