resource "consul_config_entry" "global_proxy_defaults" {
  kind = "proxy-defaults"
  name = "global"

  config_json = jsonencode({
    Config = {
      envoy_prometheus_bind_addr = "0.0.0.0:9102"
    }
  })
}

resource "consul_config_entry" "deploy_grpc_defaults" {
  kind = "service-defaults"
  name = "deploy-grpc"

  config_json = jsonencode({
    Protocol = "grpc"
  })
}

resource "consul_config_entry" "deploy_grpc_intentions" {
  kind = "service-intentions"
  name = "deploy-grpc"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "homebase-api"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/DeployService/"
            }
          },
          {
            Action = "deny"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "go_links_grpc_defaults" {
  kind = "service-defaults"
  name = "go-links-grpc"

  config_json = jsonencode({
    Protocol = "grpc"
  })
}

resource "consul_config_entry" "go_links_grpc_intentions" {
  kind = "service-intentions"
  name = "go-links-grpc"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "homebase-api"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/LinksService/"
            }
          },
          {
            Action = "deny"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "go_links_defaults" {
  kind = "service-defaults"
  name = "go-links"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "go_links_intentions" {
  kind = "service-intentions"
  name = "go-links"

  config_json = jsonencode({
    Sources = [
      {
        Name       = "ingress-http",
        Action     = "allow"
        Precedence = 9
        Type       = "consul"
      },
      {
        Name       = "*"
        Action     = "deny"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "homebase_bot_grpc_defaults" {
  kind = "service-defaults"
  name = "homebase-bot-grpc"

  config_json = jsonencode({
    Protocol = "grpc"
  })
}

resource "consul_config_entry" "homebase_bot_grpc_intentions" {
  kind = "service-intentions"
  name = "homebase-bot-grpc"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "detect-presence"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/MessagesService/"
            }
          },
          {
            Action = "deny"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "homebase_bot_defaults" {
  kind = "service-defaults"
  name = "homebase-bot"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "homebase_bot_intentions" {
  kind = "service-intentions"
  name = "homebase-bot"

  config_json = jsonencode({
    Sources = [
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "homebase_api_defaults" {
  kind = "service-defaults"
  name = "homebase-api"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "homebase_api_intentions" {
  kind = "service-intentions"
  name = "homebase-api"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "ingress-http"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathExact = "/graphql"
            }
          },
          {
            Action = "deny"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Name        = "homebase"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathExact = "/graphql"
            }
          },
          {
            Action = "deny"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "homebase_defaults" {
  kind = "service-defaults"
  name = "homebase"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "homebase_intentions" {
  kind = "service-intentions"
  name = "homebase"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "ingress-http"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "vault_proxy_defaults" {
  kind = "service-defaults"
  name = "vault-proxy"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "vault_proxy_intentions" {
  kind = "service-intentions"
  name = "vault-proxy"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "ingress-http"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}

resource "consul_config_entry" "minio_defaults" {
  kind = "service-defaults"
  name = "minio"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "minio_intentions" {
  kind = "service-intentions"
  name = "minio"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "*"
        Precedence  = 8
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
    ]
  })
}
