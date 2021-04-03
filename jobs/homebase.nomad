job "homebase" {
  datacenters = ["dc1"]

  type = "service"

  group "homebase-srv" {
    count = 2

    network {
      mode = "bridge"
    }

    service {
      name = "homebase"
      port = 3000

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "homebase-api"
              local_bind_port  = 6460
            }
          }
        }
      }
    }

    task "homebase-srv" {
      driver = "docker"

      config {
        image = "mmoriarity/homebase-srv-next"
        ports = ["http"]
      }

      resources {
        cpu    = 100
        memory = 100
      }

      env {
        GRAPHQL_URL = "http://localhost:6460/graphql"
      }
    }
  }

  group "homebase-api" {
    count = 2

    network {
      mode = "bridge"
    }

    service {
      name = "homebase-api"
      port = 6460

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
              destination_name = "go-links-grpc"
              local_bind_port  = 4241
            }
            upstreams {
              destination_name = "detect-presence-grpc"
              local_bind_port  = 2121
            }
            upstreams {
              destination_name = "deploy-grpc"
              local_bind_port  = 8481
            }
            upstreams {
              destination_name = "backup-grpc"
              local_bind_port  = 2321
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
        image   = "mmoriarity/homebase-api-srv"
        command = "/homebase-api-srv"
        args    = [
          "-prometheus-url",
          "http://10.0.0.2:9090",
        ]
      }

      resources {
        cpu    = 50
        memory = 50
      }
    }
  }

  group "homebase-bot" {
    count = 2

    network {
      mode = "bridge"
    }

    service {
      name = "homebase-bot"
      port = 6360

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
      name = "homebase-bot-grpc"
      port = 6361

      connect {
        sidecar_service {}
      }
    }

    task "homebase-bot-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/homebase-bot-srv"
        command = "/homebase-bot-srv"
        args    = [
          "-db",
          "dbname=homebase_bot host=127.0.0.1 sslmode=disable",
          "-leader-elect",
        ]
      }

      resources {
        cpu    = 50
        memory = 50
      }

      env {
        CONSUL_HTTP_ADDR = "${attr.unique.network.ip-address}:8500"
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
{{ with secret "consul/creds/homebase-bot" }}
CONSUL_HTTP_TOKEN={{ .Data.token }}
{{ end }}
EOF
        destination = "secrets/secrets.env"
        env         = true
        change_mode = "restart"
      }
    }
  }
}
