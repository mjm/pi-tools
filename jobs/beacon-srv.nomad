job "beacon" {
  datacenters = ["dc1"]

  type = "system"

  group "beacon" {
    task "beacon-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/beacon-srv"
        command = "/beacon-srv"
        args    = [
          "-proximity-uuid",
          "7298c12b-f658-445f-b1f2-5d6d582f0fb0",
          "-node-name",
          "${node.unique.name}",
        ]

        network_mode = "host"
        cap_add      = ["NET_ADMIN", "NET_RAW"]
      }

      resources {
        cpu    = 50
        memory = 40
      }
    }
  }
}
