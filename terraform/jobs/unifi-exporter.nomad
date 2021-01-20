job "unifi-exporter" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 30

  group "unifi-exporter" {
    count = 1

    network {
      port "http" {
        to = 9130
      }
    }

    service {
      name = "unifi-exporter"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "unifi-exporter" {
      driver = "docker"

      config {
        image   = "mmoriarity/unifi_exporter@${image_digests.unifi_exporter}"
        command = "/unifi_exporter"
        args    = ["-config.file=$${NOMAD_SECRETS_DIR}/config.yml"]
        ports   = ["http"]

        logging {
          type = "journald"
          config {
            tag = "unifi-exporter"
          }
        }
      }

      resources {
        cpu    = 100
        memory = 50
      }

      vault {
        policies = ["unifi-exporter"]
      }

      template {
        // language=YAML
        data        = <<EOF
listen:
  address: :9130
  metricspath: /metrics
unifi:
  address: https://10.0.0.1
  username: mmoriarity
  password: {{ with secret "kv/unifi-exporter" }}{{ .Data.data.unifi_password }}{{ end }}
  site: Default
  insecure: true
  timeout: 10s
  unifi_os: true
EOF
        destination = "secrets/config.yml"
        change_mode = "restart"
      }
    }
  }
}
