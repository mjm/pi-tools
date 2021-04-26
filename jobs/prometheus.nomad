job "prometheus" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 80

  group "alertmanager" {
    count = 1

    volume "data" {
      type      = "host"
      read_only = false
      source    = "alertmanager_data"
    }

    network {
      port "http" {
        to = 9093
      }
    }

    service {
      name = "alertmanager"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }

      check {
        type     = "http"
        path     = "/-/ready"
        interval = "15s"
        timeout  = "3s"
      }
    }

    task "alertmanager" {
      driver = "docker"

      config {
        image = "prom/alertmanager@sha256:e690a0f96fcf69c2e1161736d6bb076e22a84841e1ec8ecc87e801c70b942200"
        args  = [
          "--config.file=${NOMAD_SECRETS_DIR}/alertmanager.yml",
          "--storage.path=/alertmanager",
          "--web.external-url=https://alertmanager.homelab/",
        ]
        ports = ["http"]
      }

      resources {
        cpu    = 50
        memory = 100
      }

      volume_mount {
        volume      = "data"
        destination = "/alertmanager"
        read_only   = false
      }

      vault {
        policies = ["alertmanager"]
      }

      template {
        data            = file("prometheus/alertmanager.yml")
        destination     = "secrets/alertmanager.yml"
        left_delimiter  = "<<"
        right_delimiter = ">>"
      }

      template {
        data            = file("prometheus/pushover.tpl")
        destination     = "local/templates/pushover.tpl"
        left_delimiter  = "<<"
        right_delimiter = ">>"
      }

      template {
        data            = file("prometheus/pagerduty.tpl")
        destination     = "local/templates/pagerduty.tpl"
        left_delimiter  = "<<"
        right_delimiter = ">>"
      }
    }
  }
}
