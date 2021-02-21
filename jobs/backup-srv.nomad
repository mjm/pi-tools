job "backup-srv" {
  datacenters = ["dc1"]

  type = "service"

  group "backup" {
    count = 2

    update {
      max_parallel = 1
    }

    network {
      mode = "bridge"
    }

    service {
      name = "backup"
      port = 2320

      meta {
        metrics_path = "/metrics"
      }

      check {
        type                   = "http"
        expose                 = true
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
        sidecar_service {}
      }
    }

    task "backup-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/backup-srv"
        command = "/backup-srv"
        args    = [
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

        BORG_RSH = "ssh -o StrictHostKeyChecking=no -i ${NOMAD_SECRETS_DIR}/id_rsa"
      }

      vault {
        policies = ["borg", "tarsnap"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "kv/tarsnap" }}{{ .Data.data.key | base64Decode }}{{ end }}
EOF
        destination = "secrets/tarsnap.key"
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "kv/borg" }}{{ .Data.data.private_key }}{{ end }}
EOF
        destination = "secrets/id_rsa"
        perms       = "600"
      }
    }
  }
}
