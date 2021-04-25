job "go-links" {
  datacenters = ["dc1"]

  type = "service"

  group "go-links" {
    count = 2

    network {
      mode = "bridge"
    }

    service {
      name = "go-links"
      port = 4240

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
        sidecar_service {}
      }
    }

    service {
      name = "go-links-grpc"
      port = 4241

      connect {
        sidecar_service {}
      }
    }

    task "go-links-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/go-links-srv"
        command = "/go-links"
        args    = [
          "-db",
          "dbname=golinks host=10.0.2.102 sslmode=disable",
        ]
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["go-links"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/go-links" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
        change_mode = "restart"
      }
    }
  }
}
