job "loki" {
  datacenters = [
    "dc1",
  ]

  type = "service"

  group "loki" {
    count = 1

    network {
      mode = "bridge"
      port "http" {
        to = 3100
      }
    }

    service {
      name = "loki"
      port = 3100

      connect {
        sidecar_service {}
      }
    }

    service {
      name = "loki-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "loki" {
      driver = "docker"

      config {
        image = "grafana/loki@sha256:6afc0da6995fecf15307762d378242b65cab20d4a35b4a39397d67cad48fb7fb"
        args = [
          "-config.file=${NOMAD_TASK_DIR}/loki.yml",
        ]

        logging {
          type = "journald"
          config {
            tag = "loki"
          }
        }
      }

      resources {
        cpu = 50
        memory = 200
      }

      template {
        data = <<EOF
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    address: 127.0.0.1
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
    final_sleep: 0s
  chunk_idle_period: 5m
  chunk_retain_period: 30s
  max_transfer_retries: 0

schema_config:
  configs:
    - from: 2018-04-15
      store: boltdb
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 168h

storage_config:
  boltdb:
    directory: {{ env "NOMAD_TASK_DIR" }}/data/loki/index

  filesystem:
    directory: {{ env "NOMAD_TASK_DIR" }}/data/loki/chunks

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h

chunk_store_config:
  max_look_back_period: 0s

table_manager:
  retention_deletes_enabled: false
  retention_period: 0s
EOF
        destination = "local/loki.yml"
      }
    }
  }
}
