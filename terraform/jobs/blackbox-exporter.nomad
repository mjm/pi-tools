job "blackbox-exporter" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 70

  group "blackbox-exporter" {
    network {
      port "http" {}
    }

    service {
      name = "blackbox-exporter"
      port = "http"
    }

    task "blackbox-exporter" {
      driver = "docker"

      config {
        image = "prom/blackbox-exporter@sha256:7c3e8d34768f2db17dce800b0b602196871928977f205bbb8ab44e95a8821be5"
        args  = [
          "--config.file=${NOMAD_TASK_DIR}/blackbox.yml",
          "--web.listen-address=:${NOMAD_PORT_http}",
        ]
        ports = ["http"]

        network_mode = "host"

        logging {
          type = "journald"
          config {
            tag = "blackbox-exporter"
          }
        }
      }

      resources {
        cpu    = 100
        memory = 50
      }

      template {
        // language=YAML
        data        = <<EOF
modules:

  dns_public:
    prober: dns
    timeout: 5s
    dns:
      query_name: "www.google.com"
      query_type: "A"
      valid_rcodes:
        - NOERROR
      validate_answer_rrs:
        fail_if_not_matches_regexp:
          - "www.google.com.\t.*\tIN\tA\t.*"

  dns_private:
    prober: dns
    timeout: 5s
    dns:
      query_name: "raspberrypi.homelab"
      query_type: "A"
      valid_rcodes:
        - NOERROR
      validate_answer_rrs:
        fail_if_not_matches_regexp:
          - "raspberrypi.homelab.\t.*\tIN\tA\t10\\.0\\.0\\.2"

  dns_private_cname:
    prober: dns
    timeout: 5s
    dns:
      query_name: "homebase.homelab"
      query_type: "CNAME"
      valid_rcodes:
        - NOERROR
      validate_answer_rrs:
        fail_if_not_matches_regexp:
          - "homebase.homelab.\t.*\tIN\tCNAME\traspberrypi[2-3]?\\.homelab\\."
EOF
        destination = "local/blackbox.yml"
      }
    }
  }
}
