job "otel-collector" {
  datacenters = ["dc1"]

  type = "system"

  group "otel-collector" {
    network {
      port "healthcheck" {
        to = 13133
      }
      port "jaeger_thrift" {
        static = 14268
        to     = 14268
      }
      port "otlp_grpc" {
        static = 4317
        to     = 4317
      }
      port "otlp_grpc_2" {
        static = 55680
        to     = 55680
      }
      port "otlp_http" {
        static = 55681
        to     = 55681
      }
    }

    service {
      name = "otel-collector"
      port = "otlp_grpc"

      tags = ["grpc"]

      check {
        type                   = "http"
        path                   = "/"
        port                   = "healthcheck"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }
    }

    task "otel-collector" {
      driver = "docker"

      config {
        image   = "mmoriarity/opentelemetry-collector"
        command = "/otelcol"
        args    = [
          "--config",
          "${NOMAD_SECRETS_DIR}/config.yaml",
        ]
        ports   = ["jaeger_thrift", "otlp_grpc", "otlp_grpc_2", "otlp_http", "healthcheck"]
      }

      resources {
        cpu    = 100
        memory = 100
      }

      vault {
        policies    = ["otel-collector"]
        change_mode = "noop"
      }

      template {
        data        = file("otel-collector/otel-collector-config.yaml")
        destination = "secrets/config.yaml"
        change_mode = "restart"
      }
    }
  }
}
