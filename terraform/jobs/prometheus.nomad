job "prometheus" {
  datacenters = [
    "dc1",
  ]

  type = "service"

  group "prometheus" {
    count = 1

    volume "data" {
      type = "host"
      read_only = false
      source = "prometheus_data"
    }

    network {
      port "http" {}
    }

    service {
      name = "prometheus"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }

      check {
        type = "http"
        path = "/-/ready"
        interval = "15s"
        timeout = "3s"
      }
    }

    task "prometheus" {
      driver = "docker"

      config {
        image = "mmoriarity/prometheus@__DIGEST__"
        args = [
          "--web.listen-address=0.0.0.0:${NOMAD_PORT_http}",
          "--config.file=/etc/prometheus/prometheus.yml",
          "--storage.tsdb.path=/prometheus",
          "--web.console.libraries=/usr/share/prometheus/console_libraries",
          "--web.console.templates=/usr/share/prometheus/consoles",
          "--web.external-url=https://prometheus.homelab/",
          "--web.enable-admin-api",
        ]
        network_mode = "host"
      }

      resources {
        cpu = 200
        memory = 1500
      }

      volume_mount {
        volume      = "data"
        destination = "/prometheus"
        read_only   = false
      }
    }
  }
}
