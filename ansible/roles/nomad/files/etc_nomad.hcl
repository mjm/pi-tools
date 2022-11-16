data_dir = "/var/lib/nomad"

server {
  enabled = true
  bootstrap_expect = 3
}

acl {
  enabled = true
}

tls {
  http = true
  rpc  = true

  ca_file   = "/etc/nomad/ca.crt"
  cert_file = "/etc/nomad/agent.crt"
  key_file  = "/etc/nomad/agent.key"

  verify_server_hostname = true
  verify_https_client    = true
}

vault {
  enabled = true
  address = "http://127.0.0.1:8200"
  create_from_role = "nomad-cluster"
}

telemetry {
  prometheus_metrics = true
  publish_allocation_metrics = true
  publish_node_metrics = true
}

plugin "docker" {
  config {
    infra_image = "rancher/pause:3.2"
    allow_privileged = true
    allow_caps = ["CHOWN", "DAC_OVERRIDE", "FSETID", "FOWNER", "MKNOD", "NET_RAW", "NET_ADMIN", "SETGID", "SETUID", "SETFCAP", "SETPCAP", "NET_BIND_SERVICE", "SYS_CHROOT", "KILL", "AUDIT_WRITE"]

    volumes {
      enabled = true
    }
  }
}

client {
  enabled = true

  // deploy-srv may want to stay alive for a while, to be able to finish an in-progress deploy.
  // The default is 30s which is just not enough for that.
  max_kill_timeout = "10m"

  meta {
    "connect.sidecar_image" = "envoyproxy/envoy:v1.23.2"
  }

  host_volume "promtail_run" {
    path = "/var/lib/promtail"
    read_only = false
  }
}
