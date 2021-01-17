job "deploy" {
  datacenters = ["dc1"]

  type     = "service"
  priority = "60"

  group "deploy" {
    count = 1
    // TODO add more and introduce leader election

    update {
      max_parallel = 1
    }

    network {
      mode = "bridge"
      port "http" {
        to = 8480
      }
      port "grpc" {
        to = 8481
      }
      port "envoy_metrics_http" {
        to = 9102
      }
      port "envoy_metrics_grpc" {
        to = 9103
      }
    }

    service {
      name = "deploy"
      port = 8480

      check {
        type                   = "http"
        port                   = "http"
        path                   = "/healthz"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "jaeger-collector"
              local_bind_port  = 14268
            }
          }
        }
      }
    }

    service {
      name = "deploy-grpc"
      port = 8481

      connect {
        sidecar_service {
          proxy {
            config {
              envoy_prometheus_bind_addr = "0.0.0.0:9103"
            }
          }
        }
      }
    }

    service {
      name = "deploy-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "deploy-metrics"
      port = "envoy_metrics_http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "deploy-metrics"
      port = "envoy_metrics_grpc"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "deploy-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/deploy-srv@__DIGEST__"
        command = "/deploy-srv"
        args = [
          // TODO remove after merging this branch to main
          "-branch",
          "nomad",
        ]

        logging {
          type = "journald"
          config {
            tag = "deploy-srv"
          }
        }

        mount {
          type = "bind"
          target = "/var/run/docker.sock"
          source = "/var/run/docker.sock"
        }
      }

      resources {
        cpu    = 50
        memory = 200
      }

      env {
        CONSUL_HTTP_ADDR      = "${attr.unique.network.ip-address}:8500"
        TF_VAR_consul_address = "${attr.unique.network.ip-address}:8500"

        NOMAD_ADDR           = "http://${attr.unique.network.ip-address}:4646"
        TF_VAR_nomad_address = "http://${attr.unique.network.ip-address}:4646"

        HOSTNAME        = "${attr.unique.hostname}"
        NOMAD_CLIENT_ID = "${node.unique.id}"
      }

      vault {
        policies = ["deploy"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/deploy" }}{{ .Data.data.github_token }}{{ end }}
EOF
        destination = "secrets/github-token"
        change_mode = "restart"
      }
    }
  }
}
