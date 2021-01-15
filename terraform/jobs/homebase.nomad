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
        ports = ["http"]

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

  // TODO homebase-api
  // TODO homebase-bot
}
