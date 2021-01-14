data_dir = "/var/lib/nomad"

server {
  enabled = true
  bootstrap_expect = 3
}

client {
  enabled = true

  host_volume "postgresql_0" {
    path = "/srv/mnt/postgresql-0"
    read_only = false
  }

  meta {
    "connect.sidecar_image" = "envoyproxy/envoy:v1.16.2"
  }
}

vault {
  enabled = true
  address = "http://127.0.0.1:8200"
  create_from_role = "nomad-cluster"
}

plugin "docker" {
  config {
    infra_image = "rancher/pause:3.2"
  }
}
