job "node-exporter" {
  datacenters = [
    "dc1",
  ]

  type = "system"
  priority = 70

  group "node-exporter" {
    network {
      port "http" {
        static = 9100
        to = 9100
      }
    }

    service {
      name = "node-exporter"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "node-exporter" {
      driver = "docker"

      config {
        image = "prom/node-exporter@sha256:eb80355f0ff0a0a0f0342303cd694af28e2820d688f416049d7be7d1760a0b33"
        args = [
          "--path.procfs=/host/proc",
          "--path.sysfs=/host/sys",
          "--path.rootfs=/host/root",
          "--collector.processes",
          "--collector.systemd",
          "--collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+|var/lib/nomad/.+|run/docker/.+|snap/.+)($|/)",
          "--collector.netclass.ignored-devices=^veth",
          "--collector.netdev.device-exclude=^veth",
        ]

        logging {
          type = "journald"
          config {
            tag = "node-exporter"
          }
        }

        privileged = true
        pid_mode = "host"
        network_mode = "host"

        mount {
          type = "bind"
          target = "/host/root"
          source = "/"
          bind_options {
            propagation = "rslave" # :(
          }
        }

        mount {
          type = "bind"
          target = "/host/proc"
          source = "/proc"
          bind_options {
            propagation = "rslave" # :(
          }
        }

        mount {
          type = "bind"
          target = "/host/sys"
          source = "/sys"
          bind_options {
            propagation = "rslave" # :(
          }
        }

        mount {
          type = "bind"
          target = "/run/systemd"
          source = "/run/systemd"
          bind_options {
            propagation = "rslave" # :(
          }
        }

        mount {
          type = "bind"
          target = "/var/run/dbus"
          source = "/var/run/dbus"
        }
      }

      resources {
        cpu = 200
        memory = 50
      }
    }
  }
}
