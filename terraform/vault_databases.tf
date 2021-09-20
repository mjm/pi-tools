resource "vault_database_secret_backend_role" "grafana" {
  name    = "grafana"
  backend = "database"
  db_name = "db1"

  creation_statements = [
    "create role \"{{name}}\" with login password '{{password}}' valid until '{{expiration}}'; grant grafana_user to \"{{name}}\";"
  ]

  default_ttl = 3600
  max_ttl     = 86400
}

resource "vault_database_secret_backend_role" "paperless" {
  name    = "paperless"
  backend = "database"
  db_name = "db1"

  creation_statements = [
    "create role \"{{name}}\" with login password '{{password}}' valid until '{{expiration}}'; grant paperless_user to \"{{name}}\";"
  ]

  default_ttl = 3600
  max_ttl     = 86400
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

resource "vault_database_secret_backend_role" "phabricator" {
  name    = "phabricator"
  backend = "database"
  db_name = "db2"

  creation_statements = [
    "CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT ALL ON `phabricator\\_%`.* TO '{{name}}'@'%';"
  ]

  default_ttl = 86400
  max_ttl     = 604800
}
