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

        logging {
          type = "journald"
          config {
            tag = "prometheus"
          }
        }
      }

      resources {
        cpu = 200
        memory = 1500
      }

      volume_mount {
        volume = "data"
        destination = "/prometheus"
        read_only = false
      }
    }
  }

  group "alertmanager" {
    count = 1

    volume "data" {
      type = "host"
      read_only = false
      source = "alertmanager_data"
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
        type = "http"
        path = "/-/ready"
        interval = "15s"
        timeout = "3s"
      }
    }

    task "prometheus" {
      driver = "docker"

      config {
        image = "prom/alertmanager@sha256:e690a0f96fcf69c2e1161736d6bb076e22a84841e1ec8ecc87e801c70b942200"
        args = [
          "--config.file=${NOMAD_SECRETS_DIR}/alertmanager.yml",
          "--storage.path=/alertmanager",
          "--web.external-url=https://alertmanager.homelab/",
        ]
        ports = ["http"]

        logging {
          type = "journald"
          config {
            tag = "alertmanager"
          }
        }
      }

      resources {
        cpu = 50
        memory = 100
      }

      volume_mount {
        volume = "data"
        destination = "/alertmanager"
        read_only = false
      }

      vault {
        policies = [
          "alertmanager",
        ]
      }

      template {
        data = <<EOF
global:
  resolve_timeout: 5m

templates:
  - << env "NOMAD_TASK_DIR" >>/templates/*.tpl

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'push.phone.matt'

receivers:
  - name: 'push.phone.matt'
    pushover_configs:
<<- with secret "kv/alertmanager/pushover" >>
      - user_key: << .Data.data.user_key >>
        token: << .Data.data.token >>
<<- end >>
        title: '{{ template "pushover.title" . }}'
        priority: '{{ template "pushover.priority" . }}'

inhibit_rules: []
EOF
        destination = "secrets/alertmanager.yml"
        left_delimiter = "<<"
        right_delimiter = ">>"
      }

      template {
        data = <<EOF
{{ define "pushover.priority" }}{{ if eq .Status "firing" }}{{ if eq .CommonLabels.severity "notice" }}0{{ else }}2{{ end }}{{ else }}0{{ end }}{{ end }}

{{ define "pushover.title" }}[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonAnnotations.summary }}{{ end }}
EOF
        destination = "local/templates/pushover.tpl"
        left_delimiter = "<<"
        right_delimiter = ">>"
      }
    }
  }
}
