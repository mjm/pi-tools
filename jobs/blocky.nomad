job "blocky" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 90

  group "blocky" {
    count = 1

    network {
      port "dns" {}
      port "http" {}
    }

    service {
      name = "blocky"
      port = "dns"
      task = "blocky"

      tags = ["dns"]

      check {
        type     = "script"
        command  = "dig"
        args     = ["@${NOMAD_IP_dns}", "-p", "${NOMAD_HOST_PORT_dns}", "google.com"]
        interval = "30s"
        timeout  = "5s"
      }
    }

    service {
      name = "blocky"
      port = "http"

      tags = ["http"]

      meta {
        metrics_path = "/metrics"
      }

      check {
        type                   = "http"
        path                   = "/"
        timeout                = "5s"
        interval               = "30s"
        success_before_passing = 3
      }
    }

    task "blocky" {
      driver = "docker"

      config {
        image = "spx01/blocky@sha256:59b3661951c28db0eecd9bb2e673c798d7c861d286e7713665da862e5254c477"
        args  = [
          "/app/blocky",
          "--config",
          "${NOMAD_TASK_DIR}/config.yaml",
        ]
        ports = ["dns", "http"]
      }

      resources {
        cpu    = 200
        memory = 75
      }

      template {
        data        = file("blocky/config.yaml")
        destination = "local/config.yaml"
        change_mode = "restart"
      }
    }
  }
}
