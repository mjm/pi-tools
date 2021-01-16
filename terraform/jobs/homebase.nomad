job "homebase" {
  datacenters = ["dc1"]

  type = "service"

  group "homebase-srv" {
    count = 2

    network {
      port "http" {
        to = 8080
      }
    }

    service {
      name = "homebase"
      port = "http"
    }

    task "homebase-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/homebase-srv@__HOMEBASE_SRV_DIGEST__"
        command = "caddy"
        args    = [
          "run",
          "--config",
          "/homebase.caddy",
          "--adapter",
          "caddyfile",
        ]
        ports   = [
          "http"
        ]

        logging {
          type = "journald"
          config {
            tag = "homebase-srv"
          }
        }
      }

      resources {
        cpu    = 50
        memory = 75
      }
    }
  }

  group "homebase-api" {
    count = 2

    network {
      mode = "bridge"
      port "http" {
        to = 6460
      }
    }

    service {
      name = "homebase-api"
      port = 6460

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
              destination_name = "go-links-grpc"
              local_bind_port  = 4241
            }
            upstreams {
              destination_name = "detect-presence-grpc"
              local_bind_port  = 2121
            }
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
            }
          }
        }
      }
    }

    task "homebase-api-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/homebase-api-srv@__HOMEBASE_API_DIGEST__"
        command = "/homebase-api-srv"
      }

      resources {
        cpu    = 50
        memory = 50
      }

      env {
        HOSTNAME        = "${attr.unique.hostname}"
        NOMAD_CLIENT_ID = "${node.unique.id}"
      }
    }
  }

  group "homebase-bot" {
    count = 2

    network {
      mode = "bridge"
      port "http" {
        to = 6360
      }
      port "grpc" {
        to = 6361
      }
    }

    service {
      name = "homebase-bot"
      port = 6360

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
              destination_name = "postgresql"
              local_bind_port  = 5432
            }
            upstreams {
              destination_name = "detect-presence-grpc"
              local_bind_port  = 2121
            }
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
            }
          }
        }
      }
    }

    service {
      name = "homebase-bot-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "homebase-bot-grpc"
      port = 6361

      connect {
        sidecar_service {}
      }
    }

    task "homebase-bot-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/homebase-bot-srv@__HOMEBASE_BOT_DIGEST__"
        command = "/homebase-bot-srv"
        args    = [
          "-db",
          "dbname=homebase_bot host=localhost sslmode=disable",
          "-leader-elect",
        ]

        logging {
          type = "journald"
          config {
            tag = "homebase-bot-srv"
          }
        }
      }

      resources {
        cpu    = 50
        memory = 50
      }

      env {
        CONSUL_HTTP_ADDR = "${attr.unique.network.ip-address}:8500"
        HOSTNAME         = "${attr.unique.hostname}"
        NOMAD_CLIENT_ID  = "${node.unique.id}"
      }

      vault {
        policies = ["homebase-bot"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/homebase-bot" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
{{ with secret "kv/homebase-bot" }}
TELEGRAM_TOKEN={{ .Data.data.telegram_token | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
        change_mode = "restart"
      }
    }
  }
}
