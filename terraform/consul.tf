locals {
  consul_policies_path = "${path.module}/policies/consul"
}

resource "consul_acl_policy" "agent" {
  name        = "agent"
  description = "Allow Consul agents to make necessary internal requests"
  rules       = file("${local.consul_policies_path}/agent.hcl")
}

resource "consul_acl_policy" "dns" {
  name        = "dns"
  description = "Allow Consul DNS queries to read service and node info"
  rules       = file("${local.consul_policies_path}/dns.hcl")
}

resource "consul_acl_policy" "nomad_server" {
  name        = "nomad-server"
  description = "Allow Nomad servers to register services and create service identity tokens"
  rules       = file("${local.consul_policies_path}/nomad-server.hcl")
}

resource "consul_acl_policy" "nomad_client" {
  name        = "nomad-client"
  description = "Allow Nomad clients to register services"
  rules       = file("${local.consul_policies_path}/nomad-client.hcl")
}

resource "consul_acl_policy" "vault" {
  name        = "vault"
  description = "Allow Vault to use Consul as its backing data store"
  rules       = file("${local.consul_policies_path}/vault.hcl")
}

resource "consul_acl_policy" "prometheus" {
  name        = "prometheus"
  description = "Allow Prometheus to use Consul service discovery to find metrics endpoints"
  rules       = file("${local.consul_policies_path}/prometheus.hcl")
}

resource "consul_acl_policy" "homebase_bot" {
  name        = "homebase-bot"
  description = "Allow homebase-bot to use Consul for leader election"
  rules       = file("${local.consul_policies_path}/homebase-bot.hcl")
}

resource "consul_acl_policy" "deploy" {
  name        = "deploy"
  description = "Allow deploy-srv to use Consul for leader election"
  rules       = file("${local.consul_policies_path}/deploy.hcl")
}
