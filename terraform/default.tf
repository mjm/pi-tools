resource "nomad_job" "ingress" {
  jobspec = file("${path.module}/jobs/ingress.nomad")
}
