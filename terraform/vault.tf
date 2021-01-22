resource "vault_github_auth_backend" "github_auth" {
  organization = "mmoriarity"
}

resource "vault_github_user" "mjm" {
  backend  = vault_github_auth_backend.github_auth.id
  user     = "mjm"
  policies = ["admin"]
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

resource "vault_policy" "admin" {
  name   = "admin"
  policy = file("${path.module}/policies/admin.hcl")
}

resource "vault_policy" "alertmanager" {
  name   = "alertmanager"
  policy = file("${path.module}/policies/alertmanager.hcl")
}

resource "vault_policy" "backup" {
  name   = "backup"
  policy = file("${path.module}/policies/backup.hcl")
}

resource "vault_policy" "consul_template" {
  name   = "consul-template"
  policy = file("${path.module}/policies/consul-template.hcl")
}

resource "vault_policy" "deploy" {
  name   = "deploy"
  policy = file("${path.module}/policies/deploy.hcl")
}

resource "vault_policy" "go_links" {
  name   = "go-links"
  policy = file("${path.module}/policies/go-links.hcl")
}

resource "vault_policy" "grafana" {
  name   = "grafana"
  policy = file("${path.module}/policies/grafana.hcl")
}

resource "vault_policy" "homebase_bot" {
  name   = "homebase-bot"
  policy = file("${path.module}/policies/homebase-bot.hcl")
}

resource "vault_policy" "ingress" {
  name   = "ingress"
  policy = file("${path.module}/policies/ingress.hcl")
}

resource "vault_policy" "nomad_server" {
  name   = "nomad-server"
  policy = file("${path.module}/policies/nomad-server.hcl")
}

resource "vault_policy" "oauth_proxy" {
  name   = "oauth-proxy"
  policy = file("${path.module}/policies/oauth-proxy.hcl")
}

resource "vault_policy" "presence" {
  name   = "presence"
  policy = file("${path.module}/policies/presence.hcl")
}

resource "vault_policy" "prometheus" {
  name   = "prometheus"
  policy = file("${path.module}/policies/prometheus.hcl")
}

resource "vault_policy" "storage_readers" {
  name   = "storage-readers"
  policy = file("${path.module}/policies/storage-readers.hcl")
}

resource "vault_policy" "unifi_exporter" {
  name   = "unifi-exporter"
  policy = file("${path.module}/policies/unifi-exporter.hcl")
}
