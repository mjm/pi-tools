job "ingress" {
  datacenters = [
    "dc1",
  ]

  type = "system"

  group "ingress" {
    network {
      mode = "bridge"
      port "http" {
        to = 80
        static = 80
      }
      port "https" {
        to = 443
        static = 443
      }
    }

    service {
      name = "ingress-http"
      port = 80

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "go-links"
              local_bind_port = 4240
            }
          }
        }
      }
    }

    service {
      name = "ingress-https"
      port = 443
    }

    task "nginx" {
      driver = "docker"

      config {
        image = "nginx@sha256:763d95e3db66d9bd1bb926c029e5659ee67eb49ff57f83d331de5f5af6d2ae0c"
        volumes = [
          "local:/etc/nginx/conf.d",
          "secrets:/etc/nginx/ssl",
        ]
      }

      vault {
        policies = [
          "ingress"]
        change_mode = "noop"
      }

      template {
        data = <<EOF
upstream go-links {
  server 127.0.0.1:4240;
}

server {
    listen 80 default_server;
    server_name _;

    return 301 https://$host$request_uri;
}

server {
  listen 443 ssl;
  server_name go.homelab;

  ssl_certificate /etc/nginx/ssl/go.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/go.homelab.pem;

  location / {
    proxy_pass http://go-links;
  }
}
EOF

        destination = "local/load-balancer.conf"
        change_mode = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=go.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination = "secrets/go.homelab.pem"
        change_mode = "signal"
        change_signal = "SIGHUP"
      }
    }
  }
}