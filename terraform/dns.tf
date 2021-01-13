resource "nomad_namespace" "dns" {
  name = "dns"
}

resource "nomad_job" "named" {
  jobspec = file("${path.module}/jobs/dns/named.nomad")
}
