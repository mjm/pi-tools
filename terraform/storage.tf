resource "nomad_namespace" "storage" {
  name = "storage"
}

resource "nomad_job" "postgresql" {
  jobspec = file("${path.module}/jobs/storage/postgresql.nomad")
}
