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

    task "nginx" {
      driver = "docker"

      config {
        image = "nginx@sha256:763d95e3db66d9bd1bb926c029e5659ee67eb49ff57f83d331de5f5af6d2ae0c"
        volumes = [
          "local:/etc/nginx/conf.d",
        ]
      }

      template {
        data = <<EOF
upstream go-links {
  server 127.0.0.1:4240;
}

server {
  listen 80;
  server_name go.homelab;

  location / {
     proxy_pass http://go-links;
  }
}
EOF

        destination   = "local/load-balancer.conf"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }
  }
}
