resource "consul_acl_policy" "agent" {
  name  = "agent"
  rules = <<EOF
service_prefix "" {
  policy = "write"
}
EOF
}

resource "consul_acl_policy" "nomad_server" {
  name  = "nomad-server"
  rules = <<EOF
agent_prefix "" {
  policy = "read"
}

node_prefix "" {
  policy = "read"
}

service_prefix "" {
  policy = "write"
}

acl = "write"
EOF
}

resource "consul_acl_policy" "nomad_client" {
  name  = "nomad-client"
  rules = <<EOF
agent_prefix "" {
  policy = "read"
}

node_prefix "" {
  policy = "read"
}

service_prefix "" {
  policy = "write"
}
EOF
}

resource "consul_acl_policy" "vault" {
  name  = "vault"
  rules = <<EOF
key_prefix "vault/" {
  policy = "write"
}

node_prefix "" {
  policy = "write"
}

service "vault" {
  policy = "write"
}

agent_prefix "" {
  policy = "write"
}

session_prefix "" {
  policy = "write"
}
EOF
}

resource "consul_acl_policy" "prometheus" {
  name        = "prometheus"
  description = "Allow Prometheus to use Consul service discovery to find metrics endpoints"

  rules = <<EOF
service_prefix "" {
  policy = "read"
}

node_prefix "" {
  policy = "read"
}

agent_prefix "" {
  policy = "read"
}
EOF
}

resource "consul_acl_policy" "homebase_bot" {
  name = "homebase-bot"

  rules = <<EOF
key "service/homebase-bot/leader" {
  policy = "write"
}

session_prefix "" {
  policy = "write"
}
EOF
}

resource "consul_acl_policy" "deploy" {
  name = "deploy"

  rules = <<EOF
key "service/deploy/leader" {
  policy = "write"
}

session_prefix "" {
  policy = "write"
}
EOF
}
