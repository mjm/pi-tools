job "pushgateway" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "pushgateway" {
    count = 1

    network {
      port "http" {}
    }

    service {
      name = "pushgateway"
      port = "http"

      check {
        type                   = "http"
        path                   = "/-/ready"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }
    }

    task "pushgateway" {
      driver = "docker"

      config {
        image = "prom/pushgateway@sha256:84327d5679194898b4952009b8f407e79a82f5f39dfbdfe8959bc5b62a84af0d"
        args  = [
          "--web.listen-address=:${NOMAD_PORT_http}",
        ]
        ports = ["http"]
      }

      resources {
        cpu    = 50
        memory = 50
      }
    }
  }
}
