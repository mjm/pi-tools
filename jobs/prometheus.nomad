job "prometheus" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 80

  group "prometheus" {
    count = 1

    volume "data" {
      type      = "host"
      read_only = false
      source    = "prometheus_data"
    }

    network {
      port "http" {
        static = 9090
        to     = 9090
      }
    }

    service {
      name = "prometheus"
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

    task "prometheus" {
      driver = "docker"

      config {
        image        = "prom/prometheus@sha256:9fa25ec244e0109fdbeaff89496ac149c0539489f2f2126b9e38cf9837235be4"
        args         = [
          "--web.listen-address=:9090",
          "--config.file=${NOMAD_SECRETS_DIR}/prometheus.yml",
          "--storage.tsdb.path=/prometheus",
          "--web.console.libraries=/usr/share/prometheus/console_libraries",
          "--web.console.templates=/usr/share/prometheus/consoles",
          "--web.external-url=https://prometheus.homelab/",
          "--web.enable-admin-api",
        ]
        network_mode = "host"
      }

      resources {
        cpu    = 200
        memory = 1500
      }

      volume_mount {
        volume      = "data"
        destination = "/prometheus"
        read_only   = false
      }

      vault {
        policies      = ["prometheus"]
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = file("prometheus/prometheus.yml")
        destination   = "secrets/prometheus.yml"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data            = file("prometheus/alerts.yml")
        destination     = "local/rules/alerts.yml"
        left_delimiter  = "<<"
        right_delimiter = ">>"
      }

      template {
        data        = file("prometheus/presence.yml")
        destination = "local/rules/presence.yml"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ end }}
EOF
        destination   = "secrets/nomad.crt"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/nomad.key"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
        destination   = "secrets/nomad.ca.crt"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }
  }

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
