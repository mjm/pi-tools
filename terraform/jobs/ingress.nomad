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
            upstreams {
              destination_name = "homebase-api"
              local_bind_port = 6460
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

        logging {
          type = "journald"
          config {
            tag = "ingress"
          }
        }
      }

      vault {
        policies = [
          "ingress"
        ]
        change_mode = "noop"
      }

      template {
        data = <<EOF
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

upstream prometheus {
{{ range service "prometheus" }}
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
  server_name go.homelab;

  ssl_certificate /etc/nginx/ssl/go.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/go.homelab.pem;

  __OAUTH_LOCATIONS_SNIPPET__

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://go-links;
  }
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

  location / {
    __OAUTH_REQUEST_SNIPPET__

    proxy_pass http://homebase;
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

      template {
        data = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=homebase.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination = "secrets/homebase.homelab.pem"
        change_mode = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=prometheus.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination = "secrets/prometheus.homelab.pem"
        change_mode = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data = <<EOF
{{ with secret "pki-homelab/issue/homelab" "common_name=jaeger.homelab" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination = "secrets/jaeger.homelab.pem"
        change_mode = "signal"
        change_signal = "SIGHUP"
      }
    }
  }
}
