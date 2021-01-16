job "grafana" {
  datacenters = [
    "dc1",
  ]

  type = "service"

  group "grafana" {
    count = 3

    network {
      mode = "bridge"
      port "http" {
        to = 3000
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

    task "grafana" {
      driver = "docker"

      config {
        image = "mmoriarity/grafana@__DIGEST__"

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
        data = <<EOF
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: {{ range service "prometheus" }}http://{{ .Address }}:{{ .Port }} {{ end }}
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
    url: {{ range service "jaeger-query" }}http://{{ .Address }}:{{ .Port }} {{ end }}
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
    }
  }
}
