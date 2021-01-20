locals {
  oauth_locations_snippet = <<EOF
  location /oauth2/ {
    proxy_pass       http://oauth-proxy;
    proxy_set_header Host                    $host;
    proxy_set_header X-Real-IP               $remote_addr;
    proxy_set_header X-Scheme                $scheme;
    proxy_set_header X-Auth-Request-Redirect $scheme://$host$request_uri;
  }

  location /oauth2/auth {
    proxy_pass       http://oauth-proxy;
    proxy_set_header Host             $host;
    proxy_set_header X-Real-IP        $remote_addr;
    proxy_set_header X-Scheme         $scheme;
    # nginx auth_request includes headers but not body
    proxy_set_header Content-Length   "";
    proxy_pass_request_body           off;
  }
EOF
  oauth_request_snippet = <<EOF
    auth_request /oauth2/auth;
    error_page 401 = /oauth2/sign_in;

    # pass information via X-User and X-Email headers to backend,
    # requires running with --set-xauthrequest flag
    auth_request_set $user   $upstream_http_x_auth_request_user;
    auth_request_set $email  $upstream_http_x_auth_request_email;
    proxy_set_header X-Auth-Request-User  $user;
    proxy_set_header X-Auth-Request-Email $email;

    # if you enabled --cookie-refresh, this is needed for it to work with auth_request
    auth_request_set $auth_cookie $upstream_http_set_cookie;
    add_header Set-Cookie $auth_cookie;
EOF
}

job "ingress" {
  datacenters = [
    "dc1",
  ]

  type     = "service"
  priority = 70

  update {
    max_parallel = 1
    stagger      = "30s"
  }

  group "ingress" {
    count = 2

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

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

    proxy_pass https://nomad;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    proxy_ssl_certificate /etc/nginx/ssl/nomad.pem;
    proxy_ssl_certificate_key /etc/nginx/ssl/nomad.pem;
    proxy_ssl_trusted_certificate /etc/nginx/ssl/nomad.ca.crt;

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

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

    proxy_pass http://consul;
  }
}

server {
  listen 443 ssl;
  server_name vault.homelab;

  ssl_certificate /etc/nginx/ssl/vault.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/vault.homelab.pem;

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

    proxy_pass http://vault;
  }
}

server {
  listen 443 ssl;
  server_name go.homelab;

  ssl_certificate /etc/nginx/ssl/go.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/go.homelab.pem;
  add_header Strict-Transport-Security "max-age=2628000" always;

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

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

  ${local.oauth_locations_snippet}

  location /graphql {
    auth_request /oauth2/auth;
    error_page 401 = @graphql_fallback;

    auth_request_set $user   $upstream_http_x_auth_request_user;
    auth_request_set $email  $upstream_http_x_auth_request_email;
    proxy_set_header X-Auth-Request-User  $user;
    proxy_set_header X-Auth-Request-Email $email;

    proxy_pass http://homebase-api;
  }

  # The GraphQL API can handle receiving requests that weren't authorized, and will check for
  # the X-Auth-* headers itself to determine permissions.
  location @graphql_fallback {
    proxy_pass http://homebase-api;
  }

  location /download_app {
    ${local.oauth_request_snippet}

    proxy_pass http://detect-presence;
  }

  location / {
    ${local.oauth_request_snippet}

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

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

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

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

    proxy_pass http://prometheus;
  }
}

server {
  listen 443 ssl;
  server_name alertmanager.homelab;

  ssl_certificate /etc/nginx/ssl/alertmanager.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/alertmanager.homelab.pem;

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

    proxy_pass http://alertmanager;
  }
}

server {
  listen 443 ssl;
  server_name grafana.homelab;

  ssl_certificate /etc/nginx/ssl/grafana.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/grafana.homelab.pem;

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

    proxy_pass http://grafana;
  }
}

server {
  listen 443 ssl;
  server_name jaeger.homelab;

  ssl_certificate /etc/nginx/ssl/jaeger.homelab.pem;
  ssl_certificate_key /etc/nginx/ssl/jaeger.homelab.pem;

  ${local.oauth_locations_snippet}

  location / {
    ${local.oauth_request_snippet}

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

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/nomad.pem"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
        destination   = "secrets/nomad.ca.crt"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }
  }
}
