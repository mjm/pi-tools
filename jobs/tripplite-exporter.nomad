job "tripplite-exporter" {
  datacenters = ["dc1"]

  type     = "system"
  priority = 70

  group "tripplite-exporter" {
    count = 1

    // this needs to connect to the Tripplite UPS over USB,
    // and it's only plugged in to this machine.
    constraint {
      attribute = "${node.unique.name}"
      value     = "raspberrypi"
    }

    network {
      port "http" {
        to = 8080
      }
    }

    service {
      name = "tripplite-exporter"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "tripplite-exporter" {
      driver = "docker"

      config {
        image   = "mmoriarity/tripplite-exporter"
        command = "/tripplite_exporter"
        ports   = ["http"]

        privileged = true

        mount {
          type   = "bind"
          target = "/dev/bus/usb"
          source = "/dev/bus/usb"
        }
      }

      resources {
        cpu    = 30
        memory = 50
      }
    }
  }
}
