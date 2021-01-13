data_dir = "/var/lib/nomad"

server {
  enabled = true
  bootstrap_expect = 3
}

client {
  enabled = true
}
