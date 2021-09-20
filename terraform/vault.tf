resource "vault_github_auth_backend" "github_auth" {
  organization = "mmoriarity"
}

resource "vault_github_user" "mjm" {
  backend  = vault_github_auth_backend.github_auth.id
  user     = "mjm"
  policies = ["admin"]
}

resource "vault_auth_backend" "webauthn" {
  type = "webauthn"
}

resource "vault_generic_endpoint" "webauthn_config" {
  path = "auth/webauthn/config"

  data_json = jsonencode({
    display_name = "Matt's Homelab"
    id           = "home.mattmoriarity.com"
    origin       = "https://auth.home.mattmoriarity.com"
    token_ttl    = 259200
  })

  disable_delete       = true
  ignore_absent_fields = true
}

resource "vault_generic_endpoint" "webauthn_user_mjm" {
  path = "auth/webauthn/users/mjm"

  data_json = jsonencode({
    display_name   = "Matt Moriarity"
    token_policies = ["admin"]
  })

  ignore_absent_fields = true
}

resource "vault_mount" "kv" {
  path    = "kv"
  type    = "kv"
  options = {
    version = "2"
  }
}

resource "vault_mount" "database" {
  path = "database"
  type = "database"
}

resource "vault_mount" "pki" {
  path = "pki"
  type = "pki"

  max_lease_ttl_seconds = 315360000
}

resource "vault_mount" "pki_int" {
  path = "pki-int"
  type = "pki"

  max_lease_ttl_seconds = 157680000
}

resource "vault_pki_secret_backend_role" "nomad_cluster" {
  backend = vault_mount.pki_int.path
  name    = "nomad-cluster"

  max_ttl = 157680000

  require_cn = false
  key_usage  = ["DigitalSignature", "KeyAgreement", "KeyEncipherment"]

  allow_localhost    = true
  allow_bare_domains = true
  allow_subdomains   = true
  allow_glob_domains = false
  allowed_domains    = ["global.nomad", "nomad.service.consul"]
}

resource "vault_mount" "pki_homelab" {
  path = "pki-homelab"
  type = "pki"

  max_lease_ttl_seconds = 31536000
}

resource "vault_pki_secret_backend_config_urls" "homelab_config_urls" {
  backend = vault_mount.pki_homelab.path

  crl_distribution_points = ["http://vault.service.consul:8200/v1/${vault_mount.pki_homelab.path}/crl"]
  issuing_certificates    = ["http://vault.service.consul:8200/v1/${vault_mount.pki_homelab.path}/ca"]
}

resource "vault_pki_secret_backend_role" "homelab" {
  backend = vault_mount.pki_homelab.path
  name    = "homelab"

  max_ttl        = 604800
  generate_lease = true

  key_usage = ["DigitalSignature", "KeyAgreement", "KeyEncipherment"]

  allow_localhost    = false
  allow_bare_domains = false
  allow_subdomains   = true
  allow_glob_domains = false
  allowed_domains    = ["homelab", "home.mattmoriarity.com"]
}

resource "vault_mount" "nomad" {
  path = "nomad"
  type = "nomad"
}

resource "vault_generic_endpoint" "nomad_role_deploy" {
  path      = "nomad/role/deploy"
  data_json = jsonencode({
    policies = ["deploy"]
  })

  ignore_absent_fields = true
}

resource "vault_generic_endpoint" "nomad_lease_config" {
  path      = "nomad/config/lease"
  data_json = jsonencode({
    ttl = 86400
  })

  ignore_absent_fields = true
}

resource "vault_auth_backend" "approle" {
  type = "approle"
}

resource "vault_approle_auth_backend_role" "homelab" {
  role_name             = "homelab"
  secret_id_bound_cidrs = ["10.0.2.114/32"]
  token_policies        = ["homelab"]
}

resource "vault_database_secret_backend_role" "homelab" {
  name    = "homelab"
  backend = "database"
  db_name = "db1"

  creation_statements = [
    "create role \"{{name}}\" with login password '{{password}}' valid until '{{expiration}}'; grant homelab to \"{{name}}\";"
  ]

  default_ttl = 86400
  max_ttl     = 604800
}

resource "vault_approle_auth_backend_role" "paperless" {
  role_name             = "paperless"
  secret_id_bound_cidrs = ["10.0.2.110/32"]
  token_policies        = ["paperless"]
}

locals {
  vault_policies_path = "${path.module}/policies/vault"
}

resource "vault_policy" "admin" {
  name   = "admin"
  policy = file("${local.vault_policies_path}/admin.hcl")
}

resource "vault_policy" "alertmanager" {
  name   = "alertmanager"
  policy = file("${local.vault_policies_path}/alertmanager.hcl")
}

resource "vault_policy" "consul_template" {
  name   = "consul-template"
  policy = file("${local.vault_policies_path}/consul-template.hcl")
}

resource "vault_policy" "homelab" {
  name   = "homelab"
  policy = file("${local.vault_policies_path}/homelab.hcl")
}

resource "vault_policy" "nomad_server" {
  name   = "nomad-server"
  policy = file("${local.vault_policies_path}/nomad-server.hcl")
}

resource "vault_policy" "paperless" {
  name   = "paperless"
  policy = file("${local.vault_policies_path}/paperless.hcl")
}

resource "vault_policy" "phabricator" {
  name   = "phabricator"
  policy = file("${local.vault_policies_path}/phabricator.hcl")
}

resource "vault_policy" "prometheus" {
  name   = "prometheus"
  policy = file("${local.vault_policies_path}/prometheus.hcl")
}

resource "vault_policy" "teamcity" {
  name   = "teamcity"
  policy = file("${local.vault_policies_path}/teamcity.hcl")
}
