job "loki" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 80

  group "loki" {
    count = 1

    network {
      mode = "bridge"
    }

    service {
      name = "loki"
      port = 3100

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
        sidecar_service {}
      }
    }

    task "loki" {
      driver = "docker"

      config {
        # loki 2.2.1
        image = "grafana/loki@sha256:7d2ddbe46c11cf9778eba0abf67bc963366dcfd7bda1a123e5244187e64dafec"
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
