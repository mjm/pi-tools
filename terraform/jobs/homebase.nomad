job "homebase" {
  datacenters = [
    "dc1",
  ]

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
        image = "mmoriarity/homebase-srv@__HOMEBASE_SRV_DIGEST__"
        command = "caddy"
        args = [
          "run",
          "--config",
          "/homebase.caddy",
          "--adapter",
          "caddyfile",
        ]
        ports = [
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
        cpu = 50
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
        type = "http"
        port = "http"
        path = "/healthz"
        timeout = "3s"
        interval = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "go-links-grpc"
              local_bind_port = 4241
            }
            upstreams {
              destination_name = "detect-presence-grpc"
              local_bind_port = 2121
            }
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port = 14268
            }
          }
        }
      }
    }

    task "homebase-api-srv" {
      driver = "docker"

      config {
        image = "mmoriarity/homebase-api-srv@__HOMEBASE_API_DIGEST__"
        command = "/homebase-api-srv"
      }

      resources {
        cpu = 50
        memory = 50
      }
    }
  }

  // TODO homebase-bot
}
