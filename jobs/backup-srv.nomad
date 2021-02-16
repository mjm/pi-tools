job "backup-srv" {
  datacenters = ["dc1"]

  type     = "service"

  group "backup" {
    count = 2

    update {
      max_parallel = 1
    }

    volume "homelab_nfs" {
      type   = "host"
      source = "homelab_nfs"
    }

    network {
      mode = "bridge"
      port "http" {
        to = 2320
      }
      port "grpc" {
        to = 2321
      }
      port "envoy_metrics_http" {
        to = 9102
      }
      port "envoy_metrics_grpc" {
        to = 9103
      }
    }

    service {
      name = "backup"
      port = 2320

      check {
        type                   = "http"
        port                   = "http"
        path                   = "/healthz"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
            }
          }
        }
      }
    }

    service {
      name = "backup-grpc"
      port = 2321

      connect {
        sidecar_service {
          proxy {
            config {
              envoy_prometheus_bind_addr = "0.0.0.0:9103"
            }
          }
        }
      }
    }

    service {
      name = "backup-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "backup-metrics"
      port = "envoy_metrics_http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "backup-metrics"
      port = "envoy_metrics_grpc"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "backup-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/backup-srv"
        command = "/backup-srv"
        args = [
          "-tarsnap-keyfile",
          "${NOMAD_SECRETS_DIR}/tarsnap.key",
        ]
      }

      resources {
        cpu    = 50
        memory = 100
      }

      env {
        BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK = "yes"
      }

      volume_mount {
        volume      = "homelab_nfs"
        destination = "/backup/borg"
      }

      vault {
        policies = ["tarsnap"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/tarsnap" }}{{ .Data.data.key | base64Decode }}{{ end }}
EOF
        destination = "secrets/tarsnap.key"
      }
    }
  }
}
