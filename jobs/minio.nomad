job "minio" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 60

  group "minio" {
    count = 1

    volume "data" {
      type      = "host"
      read_only = false
      source    = "minio_0"
    }

    network {
      mode = "bridge"
    }

    service {
      name = "minio"
      port = 9000

      check {
        type     = "http"
        expose   = true
        path     = "/minio/health/cluster"
        interval = "15s"
        timeout  = "10s"
      }

      check {
        type     = "http"
        expose   = true
        path     = "/minio/health/live"
        interval = "15s"
        timeout  = "10s"

        check_restart {
          limit = 3
          grace = "120s"
        }
      }

      connect {
        sidecar_service {}
      }
    }

    task "minio" {
      driver = "docker"

      config {
        image = "minio/minio@sha256:7fe919b99b0ba1f217ce894d170816f622aea7fc32d7a2ae3765a0f0b4f95d5a"
        args  = ["server", "/data"]
      }

      resources {
        cpu    = 200
        memory = 500
      }

      volume_mount {
        volume      = "data"
        destination = "/data"
        read_only   = false
      }
    }
  }
}
