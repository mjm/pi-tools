job "unifi-exporter" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 30

  group "unifi-exporter" {
    count = 1

    network {
      port "http" {
        to = 9130
      }
    }

    service {
      name = "unifi-exporter"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "unifi-exporter" {
      driver = "docker"

      config {
        image   = "mmoriarity/unifi_exporter"
        command = "/unifi_exporter"
        args    = ["-config.file=${NOMAD_SECRETS_DIR}/config.yml"]
        ports   = ["http"]
      }

      resources {
        cpu    = 50
        memory = 30
      }

      vault {
        policies = ["unifi-exporter"]
      }

      template {
        data        = file("unifi-exporter/config.yml")
        destination = "secrets/config.yml"
        change_mode = "restart"
      }
    }
  }
}
