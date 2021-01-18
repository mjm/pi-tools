job "grafana" {
  datacenters = [
    "dc1",
  ]

  type = "service"
  priority = 70

  group "grafana" {
    count = 3

    network {
      mode = "bridge"
      port "http" {
        to = 3000
      }
      port "envoy_metrics" {
        to = 9102
      }
    }

    service {
      name = "grafana"
      port = 3000

      check {
        expose = true
        path = "/api/health"
        type = "http"
        interval = "15s"
        timeout = "3s"
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgresql"
              local_bind_port = 5432
            }
            upstreams {
              destination_name = "loki"
              local_bind_port = 3100
            }
          }
        }
      }
    }

    service {
      name = "grafana-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "grafana-metrics"
      port = "envoy_metrics"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "grafana" {
      driver = "docker"

      config {
        image = "grafana/grafana@sha256:f0817ecbf8dcf33e10cca2245bd25439433c441189bbe1ce935ac61d05f9cc6f"

        logging {
          type = "journald"
          config {
            tag = "grafana"
          }
        }
      }

      resources {
        cpu = 100
        memory = 150
      }

      env {
        GF_PATHS_CONFIG = "${NOMAD_SECRETS_DIR}/grafana.ini"
        GF_PATHS_PROVISIONING = "${NOMAD_TASK_DIR}/provisioning"
      }

      vault {
        policies = [
          "grafana"
        ]
      }

      template {
        data = <<EOF
[log]
level = info

[database]
type = postgres
host = 127.0.0.1
name = grafana
{{ with secret "database/creds/grafana" -}}
user = {{ .Data.username }}
password = """{{ .Data.password }}"""
{{ end -}}
ssl_mode = disable

[users]
auto_assign_org_role = Admin

[auth.proxy]
enabled = true
header_name = X-Auth-Request-User
auto_sign_up = true
headers = Email:X-Auth-Request-Email
enable_login_token = false
EOF
        destination = "secrets/grafana.ini"
        change_mode = "restart"
      }

      template {
        // language=YAML
        data = <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: {{ with service "prometheus" }}{{ with index . 0 }}http://{{ .Address }}:{{ .Port }}{{ end }}{{ end }}
    isDefault: true
    version: 1
    editable: false
  - name: Loki
    type: loki
    access: proxy
    url: http://127.0.0.1:3100
    version: 1
    editable: false
    jsonData:
      maxLines: 1000
  - name: Jaeger
    type: jaeger
    access: proxy
    url: {{ with service "jaeger-query" }}{{ with index . 0 }}http://{{ .Address }}:{{ .Port }}{{ end }}{{ end }}
    version: 1
    editable: false
EOF
        destination = "local/provisioning/datasources/datasources.yaml"
      }

      template {
        data = <<EOF
apiVersion: 1

providers:
  - name: dashboards
    type: file
    updateIntervalSeconds: 600
    options:
      path: {{ env "NOMAD_TASK_DIR" }}/dashboards
      foldersFromFilesStructure: true
EOF
        destination = "local/provisioning/dashboards/dashboards.yaml"
      }

      artifact {
        source = "https://raw.githubusercontent.com/mjm/pi-tools/main/monitoring/grafana/provisioning/dashboards/cluster.json"
        destination = "local/dashboards"
      }

      artifact {
        source = "https://raw.githubusercontent.com/mjm/pi-tools/main/monitoring/grafana/provisioning/dashboards/home.json"
        destination = "local/dashboards"
      }

      artifact {
        source = "https://raw.githubusercontent.com/mjm/pi-tools/main/monitoring/grafana/provisioning/dashboards/node.json"
        destination = "local/dashboards"
      }

      artifact {
        source = "https://raw.githubusercontent.com/mjm/pi-tools/main/monitoring/grafana/provisioning/dashboards/envoy.json"
        destination = "local/dashboards"
      }
    }
  }
}
