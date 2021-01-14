resource "nomad_namespace" "apps" {
  name = "apps"
}

data "docker_registry_image" "go_links_srv" {
  name = "mmoriarity/go-links-srv:latest"
}

resource "nomad_job" "go_links_srv" {
  jobspec = replace(file("${path.module}/jobs/apps/go-links-srv.nomad"), "__DIGEST__", data.docker_registry_image.go_links_srv.sha256_digest)
}
