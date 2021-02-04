storage "consul" {}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = true
}

ui = true

plugin_directory = "/usr/local/libexec/vault"

telemetry {
  disable_hostname = true
}
