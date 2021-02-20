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
    id           = "homelab"
    origin       = "https://auth.homelab"
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
  allowed_domains    = ["homelab"]
}

resource "vault_mount" "consul" {
  path = "consul"
  type = "consul"
}

resource "vault_consul_secret_backend_role" "nomad_client_server" {
  backend = vault_mount.consul.id
  name    = "nomad-client-server"

  policies = [
    consul_acl_policy.nomad_client.name,
    consul_acl_policy.nomad_server.name,
  ]
}

resource "vault_consul_secret_backend_role" "prometheus" {
  backend = vault_mount.consul.id
  name    = "prometheus"

  policies = [consul_acl_policy.prometheus.name]
}

resource "vault_consul_secret_backend_role" "backup" {
  backend = vault_mount.consul.id
  name    = "backup"

  # Backups take Consul snapshots, and saving a snapshot understandably requires a management token
  policies = ["global-management"]
}

resource "vault_consul_secret_backend_role" "homebase_bot" {
  backend = vault_mount.consul.id
  name    = "homebase-bot"

  policies = [consul_acl_policy.homebase_bot.name]
}

resource "vault_consul_secret_backend_role" "deploy" {
  backend = vault_mount.consul.id
  name    = "deploy"

  policies = [consul_acl_policy.deploy.name]
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

resource "vault_policy" "backup" {
  name   = "backup"
  policy = file("${local.vault_policies_path}/backup.hcl")
}

resource "vault_policy" "borg" {
  name   = "borg"
  policy = file("${local.vault_policies_path}/borg.hcl")
}

resource "vault_policy" "consul_exporter" {
  name   = "consul-exporter"
  policy = file("${local.vault_policies_path}/consul-exporter.hcl")
}

resource "vault_policy" "consul_template" {
  name   = "consul-template"
  policy = file("${local.vault_policies_path}/consul-template.hcl")
}

resource "vault_policy" "deploy" {
  name   = "deploy"
  policy = file("${local.vault_policies_path}/deploy.hcl")
}

resource "vault_policy" "go_links" {
  name   = "go-links"
  policy = file("${local.vault_policies_path}/go-links.hcl")
}

resource "vault_policy" "grafana" {
  name   = "grafana"
  policy = file("${local.vault_policies_path}/grafana.hcl")
}

resource "vault_policy" "homebase_bot" {
  name   = "homebase-bot"
  policy = file("${local.vault_policies_path}/homebase-bot.hcl")
}

resource "vault_policy" "ingress" {
  name   = "ingress"
  policy = file("${local.vault_policies_path}/ingress.hcl")
}

resource "vault_policy" "nomad_server" {
  name   = "nomad-server"
  policy = file("${local.vault_policies_path}/nomad-server.hcl")
}

resource "vault_policy" "vault_proxy" {
  name   = "vault-proxy"
  policy = file("${local.vault_policies_path}/vault-proxy.hcl")
}

resource "vault_policy" "presence" {
  name   = "presence"
  policy = file("${local.vault_policies_path}/presence.hcl")
}

resource "vault_policy" "prometheus" {
  name   = "prometheus"
  policy = file("${local.vault_policies_path}/prometheus.hcl")
}

resource "vault_policy" "storage_readers" {
  name   = "storage-readers"
  policy = file("${local.vault_policies_path}/storage-readers.hcl")
}

resource "vault_policy" "tarsnap" {
  name   = "tarsnap"
  policy = file("${local.vault_policies_path}/tarsnap.hcl")
}

resource "vault_policy" "unifi_exporter" {
  name   = "unifi-exporter"
  policy = file("${local.vault_policies_path}/unifi-exporter.hcl")
}
