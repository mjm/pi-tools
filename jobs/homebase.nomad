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
        image   = "mmoriarity/homebase-srv"
        command = "caddy"
        args    = [
          "run",
          "--config",
          "/homebase.caddy",
          "--adapter",
          "caddyfile",
        ]
        ports   = ["http"]
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
      port "expose" {}
      port "envoy_metrics" {
        to = 9102
      }
    }

    service {
      name = "homebase-api"
      port = 6460

      meta {
        metrics_path       = "/metrics"
        metrics_port       = "${NOMAD_HOST_PORT_expose}"
        envoy_metrics_port = "${NOMAD_HOST_PORT_envoy_metrics}"
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
            expose {
              path {
                path            = "/metrics"
                protocol        = "http"
                local_path_port = 6460
                listener_port   = "expose"
              }
            }
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
      port "expose" {}
      port "envoy_metrics_http" {
        to = 9102
      }
      port "envoy_metrics_grpc" {
        to = 9103
      }
    }

    service {
      name = "homebase-bot"
      port = 6360

      meta {
        metrics_path       = "/metrics"
        metrics_port       = "${NOMAD_HOST_PORT_expose}"
        envoy_metrics_port = "${NOMAD_HOST_PORT_envoy_metrics_http}"
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
            expose {
              path {
                path            = "/metrics"
                protocol        = "http"
                local_path_port = 6360
                listener_port   = "expose"
              }
            }
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

      meta {
        envoy_metrics_port = "${NOMAD_HOST_PORT_envoy_metrics_grpc}"
      }

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
