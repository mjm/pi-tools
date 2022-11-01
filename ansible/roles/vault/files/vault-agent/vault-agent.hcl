vault {
  address = "http://vault.service.consul:8200"
}

auto_auth {
  method "approle" {
    config = {
      role_id_file_path   = "/usr/local/etc/vault_role_id"
      secret_id_file_path = "/usr/local/etc/vault_secret_id"

      remove_secret_id_file_after_reading = false
    }
  }
}

