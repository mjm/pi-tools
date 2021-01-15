job "prometheus" {
  datacenters = [
    "dc1",
  ]

  type = "service"

  group "prometheus" {
    count = 1

    // prometheus's data only lives on this particular machine
//    constraint {
//      attribute = "${node.unique.name}"
//      value = "raspberrypi"
//    }

    volume "data" {
      type = "host"
      read_only = false
      source = "prometheus_data"
    }

    network {
      port "http" {
        to = 9090
      }
    }

    service {
      name = "prometheus"
      port = "http"

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
          "--config.file=/etc/prometheus/prometheus.yml",
          "--storage.tsdb.path=/prometheus",
          "--web.console.libraries=/usr/share/prometheus/console_libraries",
          "--web.console.templates=/usr/share/prometheus/consoles",
          "--web.external-url=https://prometheus.homelab/",
          "--web.enable-admin-api",
        ]
        ports = ["http"]
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
