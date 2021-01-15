resource "nomad_job" "named" {
  jobspec = file("${path.module}/jobs/named.nomad")
}

resource "nomad_job" "postgresql" {
  jobspec = file("${path.module}/jobs/postgresql.nomad")
}

resource "nomad_job" "go_links_srv" {
  jobspec = replace(file("${path.module}/jobs/go-links-srv.nomad"), "__DIGEST__", data.docker_registry_image.go_links_srv.sha256_digest)
}

resource "nomad_job" "ingress" {
  jobspec = file("${path.module}/jobs/ingress.nomad")
}

