job "consul-exporter" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "consul-exporter" {
    network {
      port "http" {}
    }

    service {
      name = "consul-exporter"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    task "consul-exporter" {
      driver = "docker"

      config {
        image = "prom/consul-exporter@sha256:4e45d018f2fd35afbc3c0c79aa6fe9f43642f9fe49170aca989998015c76c922"
        args  = [
          "--web.listen-address=:${NOMAD_PORT_http}",
          "--consul.server=${attr.unique.network.ip-address}:8500",
        ]
        ports = ["http"]
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["consul-exporter"]
      }

      template {
        data        = <<EOF
{{ with secret "consul/creds/prometheus" }}
CONSUL_HTTP_TOKEN={{ .Data.token }}
{{ end }}
EOF
        destination = "secrets/consul.env"
        env         = true
        change_mode = "restart"
      }
    }
  }
}
