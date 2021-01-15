// TODO move these into the job spec as soon as that's possible
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
EOF
}

resource "nomad_job" "named" {
  jobspec = file("${path.module}/jobs/named.nomad")
}

resource "nomad_job" "prometheus" {
  jobspec = replace(file("${path.module}/jobs/prometheus.nomad"), "__DIGEST__", data.docker_registry_image.prometheus.sha256_digest)
}

resource "nomad_job" "node_exporter" {
  jobspec = file("${path.module}/jobs/node-exporter.nomad")
}

resource "nomad_job" "jaeger" {
  jobspec = file("${path.module}/jobs/jaeger.nomad")
}

resource "nomad_job" "postgresql" {
  jobspec = file("${path.module}/jobs/postgresql.nomad")
}

resource "nomad_job" "beacon_srv" {
  jobspec = replace(file("${path.module}/jobs/beacon-srv.nomad"), "__DIGEST__", data.docker_registry_image.beacon_srv.sha256_digest)
}

resource "nomad_job" "go_links_srv" {
  jobspec = replace(file("${path.module}/jobs/go-links-srv.nomad"), "__DIGEST__", data.docker_registry_image.go_links_srv.sha256_digest)
}

resource "nomad_job" "homebase" {
  jobspec = replace(file("${path.module}/jobs/homebase.nomad"), "__HOMEBASE_SRV_DIGEST__", data.docker_registry_image.homebase_srv.sha256_digest)
}

resource "nomad_job" "oauth_proxy" {
  jobspec = file("${path.module}/jobs/oauth-proxy.nomad")
}

resource "nomad_job" "ingress" {
  jobspec = replace(replace(file("${path.module}/jobs/ingress.nomad"), "__OAUTH_LOCATIONS_SNIPPET__", local.oauth_locations_snippet),
  "__OAUTH_REQUEST_SNIPPET__", local.oauth_request_snippet)
}

