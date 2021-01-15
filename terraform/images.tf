data "docker_registry_image" "beacon_srv" {
  name = "mmoriarity/beacon-srv:latest"
}

data "docker_registry_image" "go_links_srv" {
  name = "mmoriarity/go-links-srv:latest"
}

data "docker_registry_image" "homebase_srv" {
  name = "mmoriarity/homebase-srv:latest"
}

data "docker_registry_image" "prometheus" {
  name = "mmoriarity/prometheus:latest"
}
