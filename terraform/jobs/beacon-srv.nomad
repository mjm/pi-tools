job "beacon" {
  datacenters = ["dc1"]

  type = "system"

  group "beacon" {
    task "beacon-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/beacon-srv@__DIGEST__"
        command = "/beacon-srv"
        args    = [
          "-proximity-uuid",
          "7298c12b-f658-445f-b1f2-5d6d582f0fb0",
          "-node-name",
          "${node.unique.name}",
        ]

        network_mode = "host"
        cap_add      = ["NET_ADMIN", "NET_RAW"]

        logging {
          type = "journald"
          config {
            tag = "beacon-srv"
          }
        }
      }

      resources {
        cpu    = 50
        memory = 40
      }

      env {
        HOSTNAME        = "${attr.unique.hostname}"
        NOMAD_CLIENT_ID = "${node.unique.id}"
      }
    }
  }
}
