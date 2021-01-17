job "pihole" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 90

  group "pihole" {
    count = 1

    network {
      port "dns" {
        to = 53
      }
      port "http" {
        to = 80
      }
      port "https" {
        to = 443
      }
    }

    service {
      name = "pihole"
      port = "dns"

      tags = ["dns"]
    }

    service {
      name = "pihole"
      port = "http"

      tags = ["http"]
    }

    service {
      name = "pihole"
      port = "https"

      tags = ["https"]
    }

    task "pihole" {
      driver = "docker"

      config {
        image = "pihole/pihole@sha256:d899e498cb7beb5784aef9d62e1bb7f28b4ec3306cf96f4d1bc6cd62e26bf599"
        ports = ["dns", "http", "https"]

        logging {
          type = "journald"
          config {
            tag = "pihole"
          }
        }

        mount {
          type = "bind"
          target = "/etc/dnsmasq.d"
          source = "/srv/mnt/pihole-data/dnsmasq.d"
        }

        mount {
          type = "bind"
          target = "/etc/pihole"
          source = "/srv/mnt/pihole-data/pihole"
        }
      }

      env {
        ADMIN_EMAIL               = "matt@mattmoriarity.com"
        TZ                        = "America/Denver"
        CONDITIONAL_FORWARDING    = "true"
        CONDITIONAL_FORWARDING_IP = "10.0.0.1"
        VIRTUAL_HOST              = "pihole.homelab"
      }

      resources {
        cpu    = 50
        memory = 200
      }
    }
  }
}
