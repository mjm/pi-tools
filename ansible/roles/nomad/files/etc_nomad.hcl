data_dir = "/var/lib/nomad"

server {
  enabled = true
  bootstrap_expect = 3
}

client {
  enabled = true
}

vault {
  enabled = true
  address = "http://127.0.0.1:8200"
  create_from_role = "nomad-cluster"
}
