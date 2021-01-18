job "ingress" {
  datacenters = [
    "dc1",
  ]

  type     = "system"
  priority = 70

  update {
    max_parallel = 1
    stagger      = "30s"
  }

  group "ingress" {
    network {
      mode = "bridge"
      port "http" {
        to     = 80
        static = 80
      }
      port "https" {
        to     = 443
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
              destination_name = "detect-presence"
              local_bind_port  = 2120
            }
            upstreams {
              destination_name = "detect-presence-grpc"
              local_bind_port  = 2121
            }
            upstreams {
              destination_name = "go-links"
              local_bind_port  = 4240
            }
            upstreams {
              destination_name = "homebase-api"
              local_bind_port  = 6460
            }
            upstreams {
              destination_name = "grafana"
              local_bind_port  = 3000
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
        image   = "nginx@sha256:763d95e3db66d9bd1bb926c029e5659ee67eb49ff57f83d331de5f5af6d2ae0c"
        volumes = [
          "local:/etc/nginx/conf.d",
          "secrets:/etc/nginx/ssl",
        ]

        logging {
          type = "journald"
          config {
            tag = "ingress"
          }
        }
      }

      vault {
        policies    = [
          "ingress"
        ]
        change_mode = "noop"
      }

      template {
        data          = <<EOF
# https://github.com/envoyproxy/envoy/issues/2506#issuecomment-362558239
proxy_http_version 1.1;

upstream nomad {
  ip_hash;
{{ range service "http.nomad" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream consul {
{{ range service "consul" }}
  server {{ .Address }}:8500;
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream vault {
{{ range service "vault" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream go-links {
  server 127.0.0.1:4240;
}

upstream homebase {
{{ range service "homebase" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream homebase-api {
  server 127.0.0.1:6460;
}

upstream detect-presence {
  server 127.0.0.1:2120;
}

upstream detect-presence-grpc {
  server 127.0.0.1:2121;
}

upstream pihole {
{{ range service "http.pihole" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream grafana {
  server 127.0.0.1:3000;
}

upstream prometheus {
{{ range service "prometheus" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream alertmanager {
{{ range service "alertmanager" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream jaeger-query {
{{ range service "jaeger-query" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

upstream oauth-proxy {
{{ range service "oauth-proxy" }}
  server {{ .Address }}:{{ .Port }};
{{ else }}server 127.0.0.1:65535; # force a 502
{{ end }}
}

server {
    listen 80 default_server;
    server_name _;

    return 301 https://$host$request_uri;
}

server {
  listen 443 ssl;
  server_name nomad.homelab;

  ssl_certificate /etc/nginx/ssl/nomad.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/nomad.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://nomad;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    # Nomad blocking queries will remain open for a default of 5 minutes.
    # Increase the proxy timeout to accommodate this timeout with an
    # additional grace period.
    proxy_read_timeout 310s;

    # Nomad log streaming uses streaming HTTP requests. In order to
    # synchronously stream logs from Nomad to NGINX to the browser
    # proxy buffering needs to be turned off.
    proxy_buffering off;

    # The Upgrade and Connection headers are used to establish
    # a WebSockets connection.
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";

    # The default Origin header will be the proxy address, which
    # will be rejected by Nomad. It must be rewritten to be the
    # host address instead.
    proxy_set_header Origin "${scheme}://${proxy_host}";
  }
}

server {
  listen 443 ssl;
  server_name consul.homelab;

  ssl_certificate /etc/nginx/ssl/consul.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/consul.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://consul;
  }
}

server {
  listen 443 ssl;
  server_name vault.homelab;

  ssl_certificate /etc/nginx/ssl/vault.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/vault.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://vault;
  }
}

server {
  listen 443 ssl;
  server_name go.homelab;

  ssl_certificate /etc/nginx/ssl/go.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/go.homelab.pem;
  add_header Strict-Transport-Security "max-age=2628000" always;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://go-links;
  }
}

server {
  listen 80;
  server_name go;

  return 301 https://go.homelab$request_uri;
}

server {
  listen 443 ssl default_server;
  server_name homebase.homelab;

  ssl_certificate /etc/nginx/ssl/homebase.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/homebase.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location /graphql {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://homebase-api;
  }

  location /download_app {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://detect-presence;
  }

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://homebase;
  }
}

server {
  listen 443 ssl http2;
  server_name detect-presence-grpc.homelab;

  ssl_certificate /etc/nginx/ssl/detect-presence-grpc.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/detect-presence-grpc.homelab.pem;

  location / {
    grpc_pass grpc://detect-presence-grpc;
  }
}

server {
  listen 443 ssl;
  server_name pihole.homelab;

  ssl_certificate /etc/nginx/ssl/pihole.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/pihole.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    # Needed to have pi-hole present the admin UI instead of a "you've been blocked" page
    proxy_set_header Host $host;
    proxy_pass http://pihole;
  }
}

server {
  listen 443 ssl;
  server_name prometheus.homelab;

  ssl_certificate /etc/nginx/ssl/prometheus.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/prometheus.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://prometheus;
  }
}

server {
  listen 443 ssl;
  server_name alertmanager.homelab;

  ssl_certificate /etc/nginx/ssl/alertmanager.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/alertmanager.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://alertmanager;
  }
}

server {
  listen 443 ssl;
  server_name grafana.homelab;

  ssl_certificate /etc/nginx/ssl/grafana.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/grafana.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://grafana;
  }
}

server {
  listen 443 ssl;
  server_name jaeger.homelab;

  ssl_certificate /etc/nginx/ssl/jaeger.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/jaeger.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://jaeger-query;
  }
}
EOF
        destination   = "local/load-balancer.conf"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=nomad.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/nomad.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=consul.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/consul.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=vault.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/vault.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=go.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/go.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=homebase.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/homebase.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=detect-presence-grpc.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/detect-presence-grpc.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=pihole.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/pihole.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=grafana.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/grafana.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=prometheus.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/prometheus.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=alertmanager.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/alertmanager.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=jaeger.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/jaeger.homelab.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }
  }
}
