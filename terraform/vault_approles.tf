resource "vault_approle_auth_backend_role" "guacamole" {
  backend = vault_auth_backend.approle.id

  role_name      = "guacamole"
  token_policies = ["guacamole"]
}

resource "vault_approle_auth_backend_role" "homelab" {
  backend = vault_auth_backend.approle.id

  role_name      = "homelab"
  token_policies = ["homelab"]
}

resource "vault_approle_auth_backend_role" "paperless" {
  backend = vault_auth_backend.approle.id

  role_name      = "paperless"
  token_policies = ["paperless"]
}

resource "vault_approle_auth_backend_role" "phabricator" {
  backend = vault_auth_backend.approle.id

  role_name      = "phabricator"
  token_policies = ["phabricator"]
}

resource "vault_approle_auth_backend_role" "teamcity" {
  backend = vault_auth_backend.approle.id

  role_name      = "teamcity"
  token_policies = ["teamcity"]
}

resource "vault_approle_auth_backend_role" "prometheus" {
  backend = vault_auth_backend.approle.id

  role_name      = "prometheus"
  token_policies = ["prometheus", "alertmanager"]
}
