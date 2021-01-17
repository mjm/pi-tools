resource "consul_config_entry" "global_proxy_defaults" {
  kind = "proxy-defaults"
  name = "global"

  config_json = jsonencode({
    Config = {
      envoy_prometheus_bind_addr = "0.0.0.0:9102"
    }
  })
}

resource "consul_config_entry" "detect_presence_grpc_defaults" {
  kind = "service-defaults"
  name = "detect-presence-grpc"

  config_json = jsonencode({
    Protocol = "grpc"
  })
}

resource "consul_config_entry" "detect_presence_grpc_intentions" {
  kind = "service-intentions"
  name = "detect-presence-grpc"

  config_json = jsonencode({
    Sources = [
      {
        Name   = "*"
        Action = "deny"
      },
      {
        Name        = "ingress-http"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathExact = "/TripsService/RecordTrips"
            }
          },
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/grpc.reflection.v1alpha.ServerReflection/"
            }
          },
          {
            Action = "deny"
            HTTP   = {
              PathPrefix = "/"
            }
          }
        ]
      },
      {
        Name        = "homebase-bot"
        Permissions = [
          {
            Action = "deny"
            HTTP   = {
              PathExact = "/TripsService/RecordTrips"
            }
          },
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/TripsService/"
            }
          },
        ],
      },
      {
        Name        = "homebase-api"
        Permissions = [
          {
            Action = "deny"
            HTTP   = {
              PathExact = "/TripsService/RecordTrips"
            }
          },
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/TripsService/"
            }
          },
        ],
      },
    ]
  })
}

resource "consul_config_entry" "detect_presence_defaults" {
  kind = "service-defaults"
  name = "detect-presence"

  config_json = jsonencode({
    Protocol = "http"
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
        Action = "deny"
        Name   = "*"
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
        Name   = "ingress-http",
        Action = "allow"
      },
      {
        Name   = "*"
        Action = "deny"
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
        Action = "deny"
        Name   = "*"
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
        Action = "deny"
        Name   = "*"
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
        Action = "deny"
        Name   = "*"
      },
    ]
  })
}

resource "consul_config_entry" "ingress_http_defaults" {
  kind = "service-defaults"
  name = "ingress-http"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "ingress_http_intentions" {
  kind = "service-intentions"
  name = "ingress-http"

  // This service is only there as an identity for making upstream requests.
  // Nothing should be making requests into it.
  config_json = jsonencode({
    Sources = [
      {
        Action = "deny"
        Name   = "*"
      }
    ]
  })
}

resource "consul_config_entry" "postgresql_defaults" {
  kind = "service-defaults"
  name = "postgresql"

  config_json = jsonencode({
    Protocol = "tcp"
  })
}

resource "consul_config_entry" "postgresql_intentions" {
  kind = "service-intentions"
  name = "postgresql"

  config_json = jsonencode({
    Sources = [
      {
        Action = "allow"
        Name   = "detect-presence"
      },
      {
        Action = "allow"
        Name   = "go-links"
      },
      {
        Action = "allow"
        Name   = "grafana"
      },
      {
        Action = "allow"
        Name   = "homebase-bot"
      },
      {
        Action = "deny"
        Name   = "*"
      },
    ]
  })
}

resource "consul_config_entry" "jaeger_collector_defaults" {
  kind = "service-defaults"
  name = "jaeger-collector"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "jaeger_collector_intentions" {
  kind = "service-intentions"
  name = "jaeger-collector"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "*"
        Description = "Allow any service to send traces to the collector"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathExact = "/api/traces"
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
    ]
  })
}
