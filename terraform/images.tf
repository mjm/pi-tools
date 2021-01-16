data "docker_registry_image" "beacon_srv" {
  name = "mmoriarity/beacon-srv:latest"
}

data "docker_registry_image" "detect_presence_srv" {
  name = "mmoriarity/detect-presence-srv:latest"
}

data "docker_registry_image" "go_links_srv" {
  name = "mmoriarity/go-links-srv:latest"
}

data "docker_registry_image" "grafana" {
  name = "mmoriarity/grafana:latest"
}

data "docker_registry_image" "homebase_srv" {
  name = "mmoriarity/homebase-srv:latest"
}

data "docker_registry_image" "homebase_api_srv" {
  name = "mmoriarity/homebase-api-srv:latest"
}

data "docker_registry_image" "homebase_bot_srv" {
  name = "mmoriarity/homebase-bot-srv:latest"
}
