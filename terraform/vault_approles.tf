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
