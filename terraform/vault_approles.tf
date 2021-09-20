resource "vault_approle_auth_backend_role" "homelab" {
  backend = vault_auth_backend.approle.id

  role_name             = "homelab"
  secret_id_bound_cidrs = ["10.0.2.114/32"]
  token_policies        = ["homelab"]
}

resource "vault_approle_auth_backend_role" "paperless" {
  backend = vault_auth_backend.approle.id

  role_name             = "paperless"
  secret_id_bound_cidrs = ["10.0.2.110/32"]
  token_policies        = ["paperless"]
}

resource "vault_approle_auth_backend_role" "phabricator" {
  backend = vault_auth_backend.approle.id

  role_name             = "phabricator"
  secret_id_bound_cidrs = ["10.0.2.111/32"]
  token_policies        = ["phabricator"]
}

resource "vault_approle_auth_backend_role" "teamcity" {
  backend = vault_auth_backend.approle.id

  role_name             = "teamcity"
  secret_id_bound_cidrs = ["10.0.2.113/32"]
  token_policies        = ["teamcity"]
}
