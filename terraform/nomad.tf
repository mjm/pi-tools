locals {
  nomad_policies_path = "${path.module}/policies/nomad"
}

resource "nomad_acl_policy" "anonymous" {
  name        = "anonymous"
  description = "Anonymous policy with read-only access"
  rules_hcl   = file("${local.nomad_policies_path}/anonymous.hcl")
}

resource "nomad_acl_policy" "deploy" {
  name        = "deploy"
  description = "Allow deploy-srv to submit jobs"
  rules_hcl   = file("${local.nomad_policies_path}/deploy.hcl")
}
