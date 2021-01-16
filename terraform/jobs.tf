// TODO move these into the job spec as soon as that's possible
locals {
  oauth_locations_snippet = file("${path.module}/templates/oauth_locations.conf")
  oauth_request_snippet = file("${path.module}/templates/oauth_request.conf")
}

resource "nomad_job" "named" {
  jobspec = file("${path.module}/jobs/named.nomad")
}

resource "nomad_job" "prometheus" {
  jobspec = replace(file("${path.module}/jobs/prometheus.nomad"), "__DIGEST__", data.docker_registry_image.prometheus.sha256_digest)
}

resource "nomad_job" "loki" {
  jobspec = file("${path.module}/jobs/loki.nomad")
}

resource "nomad_job" "promtail" {
  jobspec = file("${path.module}/jobs/promtail.nomad")
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

resource "nomad_job" "grafana" {
  jobspec = replace(file("${path.module}/jobs/grafana.nomad"), "__DIGEST__", data.docker_registry_image.grafana.sha256_digest)
}

resource "nomad_job" "beacon_srv" {
  jobspec = replace(file("${path.module}/jobs/beacon-srv.nomad"), "__DIGEST__", data.docker_registry_image.beacon_srv.sha256_digest)
}

resource "nomad_job" "go_links_srv" {
  jobspec = replace(file("${path.module}/jobs/go-links-srv.nomad"), "__DIGEST__", data.docker_registry_image.go_links_srv.sha256_digest)
}

resource "nomad_job" "homebase" {
  jobspec = replace(replace(file("${path.module}/jobs/homebase.nomad"), "__HOMEBASE_SRV_DIGEST__", data.docker_registry_image.homebase_srv.sha256_digest),
  "__HOMEBASE_API_DIGEST__", data.docker_registry_image.homebase_api_srv.sha256_digest)
}

resource "nomad_job" "oauth_proxy" {
  jobspec = file("${path.module}/jobs/oauth-proxy.nomad")
}

resource "nomad_job" "ingress" {
  jobspec = replace(replace(file("${path.module}/jobs/ingress.nomad"), "__OAUTH_LOCATIONS_SNIPPET__", local.oauth_locations_snippet),
  "__OAUTH_REQUEST_SNIPPET__", local.oauth_request_snippet)
}

