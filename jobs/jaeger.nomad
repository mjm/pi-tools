job "jaeger" {
  datacenters = ["dc1"]

  type = "service"

  group "jaeger" {
    count = 1

    network {
      mode = "bridge"
      port "admin_http" {
        to = 14269
      }
      port "collector_http" {
        to = 14268
      }
      port "query_http" {
        to = 16686
      }
      port "envoy_metrics_collector" {
        to = 9102
      }
    }

    service {
      name = "jaeger-admin"
      port = "admin_http"

      meta {
        metrics_path = "/metrics"
      }

      check {
        type                   = "http"
        path                   = "/"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }
    }

    service {
      name = "jaeger-collector"
      port = 14268

      // use connect for the collector to make it easier to connect to from services
      connect {
        sidecar_service {}
      }
    }

    service {
      name = "jaeger-collector-metrics"
      port = "envoy_metrics_collector"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "jaeger-query"
      port = "query_http"
    }

    task "jaeger" {
      driver = "docker"

      config {
        image = "querycapistio/all-in-one@sha256:ad4552a9facb5e71ea2ca296fb92cf510e97783ad5068f5d23a6b169edb4a9dd"
        ports = ["admin-http", "query-http"]
      }
    }
  }
}
