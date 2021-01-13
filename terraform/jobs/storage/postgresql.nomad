job "postgresql" {
  datacenters = [
    "dc1"
  ]
  namespace = "storage"

  type = "service"

  group "postgresql" {
    count = 1

    volume "data" {
      type = "host"
      read_only = false
      source = "postgresql_0"
    }

    restart {
      attempts = 30
      interval = "15m"
      delay    = "25s"
      mode     = "delay"
    }

    network {
      port "db" {
        to = 5432
      }
    }

    service {
      name = "postgresql"
      port = "db"

      check {
        type     = "tcp"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "postgresql" {
      driver = "docker"

      config {
        image = "postgres@sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"
        ports = ["db"]
      }

      env {
        POSTGRES_PASSWORD_FILE = "${NOMAD_SECRETS_DIR}/pg-password.txt"
      }

      resources {
        cpu = 100
        memory = 500
      }

      volume_mount {
        volume      = "data"
        destination = "/var/lib/postgresql/data"
        read_only   = false
      }

      template {
        data = "{{ with secret \"kv/storage/postgresql\" }}{{ index .Data.data \"pg-password\" }}{{ end }}"
        destination = "secrets/pg-password.txt"
      }

      vault {
        policies = ["storage-readers"]
        change_mode = "noop"
      }
    }
  }
}
