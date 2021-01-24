job "tarsnap-delete" {
  datacenters = ["dc1"]

  type = "batch"

  priority = 70

  parameterized {
    payload = "required"
  }

  meta {
    logging_tag = "backup-tarsnap"
  }

  group "backup" {
    count = 1

    constraint {
      attribute = "${node.unique.name}"
      value     = "raspberrypi"
    }

    task "delete-backup" {
      driver = "docker"

      config {
        image    = "mmoriarity/tarsnap@sha256:4deeb35783541c160a09cb7a58489a7bf57bb456f4efab83e0cbd663a60bbf50"
        command  = "sh"
        args     = [
          "-c",
          "cat ${NOMAD_TASK_DIR}/archives.txt | xargs -n1 tarsnap --keyfile ${NOMAD_SECRETS_DIR}/tarsnap.key --cachedir /var/lib/tarsnap/cache --no-default-config -v -d -f",
        ]

        mount {
          type   = "bind"
          target = "/var/lib/tarsnap/cache"
          source = "/var/lib/tarsnap/cache"
        }
      }

      resources {
        cpu    = 100
        memory = 100
      }

      dispatch_payload {
        file = "archives.txt"
      }

      vault {
        policies = ["backup"]
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
