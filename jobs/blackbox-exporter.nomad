job "blackbox-exporter" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "blackbox-exporter" {
    network {
      port "http" {}
    }

    service {
      name = "blackbox-exporter"
      port = "http"
    }

    task "blackbox-exporter" {
      driver = "docker"

      config {
        image = "prom/blackbox-exporter@sha256:7c3e8d34768f2db17dce800b0b602196871928977f205bbb8ab44e95a8821be5"
        args  = [
          "--config.file=${NOMAD_TASK_DIR}/blackbox.yml",
          "--web.listen-address=:${NOMAD_PORT_http}",
        ]
        ports = ["http"]

        network_mode = "host"
      }

      resources {
        cpu    = 100
        memory = 50
      }

      template {
        data        = file("blackbox-exporter/blackbox.yml")
        destination = "local/blackbox.yml"
      }
    }
  }
}
