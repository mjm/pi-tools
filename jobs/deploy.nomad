job "deploy" {
  datacenters = ["dc1"]

  type     = "service"
  priority = "60"

  group "deploy" {
    count = 2

    update {
      max_parallel = 1
    }

    network {
      mode = "bridge"
    }

    service {
      name = "deploy"
      port = 8480

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
        sidecar_service {}
      }
    }

    task "deploy-srv" {
      driver = "docker"

      config {
        image   = "mmoriarity/deploy-srv"
        command = "/deploy-srv"
        args    = [
          "-leader-elect",
        ]
      }

      resources {
        cpu    = 50
        memory = 100
      }

      env {
        CONSUL_HTTP_ADDR  = "${attr.unique.network.ip-address}:8500"
        NOMAD_ADDR        = "https://nomad.service.consul:4646"
        NOMAD_CACERT      = "${NOMAD_SECRETS_DIR}/nomad.ca.crt"
        NOMAD_CLIENT_CERT = "${NOMAD_SECRETS_DIR}/nomad.crt"
        NOMAD_CLIENT_KEY  = "${NOMAD_SECRETS_DIR}/nomad.key"
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

      template {
        data        = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ end }}
EOF
        destination = "secrets/nomad.crt"
      }

      template {
        data        = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.private_key }}
{{ end }}
EOF
        destination = "secrets/nomad.key"
      }

      template {
        data        = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
        destination = "secrets/nomad.ca.crt"
      }

      template {
        data        = <<EOF
{{ with secret "consul/creds/deploy" }}
CONSUL_HTTP_TOKEN={{ .Data.token }}
{{ end }}
{{ with secret "nomad/creds/deploy" }}
NOMAD_TOKEN={{ .Data.secret_id }}
{{ end }}
{{ with secret "kv/pushover" }}
PUSHOVER_USER_KEY={{ .Data.data.user_key }}
PUSHOVER_TOKEN={{ .Data.data.token }}
{{ end }}
EOF
        destination = "secrets/deploy.env"
        env         = true
        change_mode = "restart"
      }
    }
  }
}
