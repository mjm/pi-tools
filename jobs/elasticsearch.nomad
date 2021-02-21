job "elasticsearch" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 60

  group "elasticsearch" {
    count = 1

    volume "data" {
      type      = "host"
      read_only = false
      source    = "elasticsearch_0"
    }

    network {
      mode = "bridge"
    }

    service {
      name = "elasticsearch"
      port = 9200

      check {
        type     = "http"
        expose   = true
        path     = "/_cluster/health"
        interval = "15s"
        timeout  = "3s"
      }

      connect {
        sidecar_service {}
      }
    }

    task "elasticsearch" {
      driver = "docker"

      config {
        image = "docker.elastic.co/elasticsearch/elasticsearch@sha256:379f51333f227286d4db5a3dfc72bb88aa55a36199c8bc825536a838c00090ac"
      }

      env {
        ES_JAVA_OPTS = "-Xms512m -Xmx512m"
      }

      resources {
        cpu    = 200
        memory = 1500
      }

      volume_mount {
        volume      = "data"
        destination = "/usr/share/elasticsearch/data"
        read_only   = false
      }

      template {
        data        = <<EOF
discovery.type=single-node
EOF
        destination = "local/elastic.env"
        env         = true
      }
    }
  }
}
