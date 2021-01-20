// These datasources let the Nomad job templates specify the latest digest
// for the images that we build ourselves.

locals {
  image_digests = {
    beacon_srv          = data.docker_registry_image.beacon_srv.sha256_digest
    deploy_srv          = data.docker_registry_image.deploy_srv.sha256_digest
    detect_presence_srv = data.docker_registry_image.detect_presence_srv.sha256_digest
    go_links_srv        = data.docker_registry_image.go_links_srv.sha256_digest
    homebase_api        = data.docker_registry_image.homebase_api_srv.sha256_digest
    homebase_bot        = data.docker_registry_image.homebase_bot_srv.sha256_digest
    homebase_srv        = data.docker_registry_image.homebase_srv.sha256_digest
    prometheus_backup   = data.docker_registry_image.prometheus_backup.sha256_digest
    tripplite_exporter  = data.docker_registry_image.tripplite_exporter.sha256_digest
    unifi_exporter      = data.docker_registry_image.unifi_exporter.sha256_digest
  }
}

data "docker_registry_image" "beacon_srv" {
  name = "mmoriarity/beacon-srv:latest"
}

data "docker_registry_image" "deploy_srv" {
  name = "mmoriarity/deploy-srv:latest"
}

data "docker_registry_image" "detect_presence_srv" {
  name = "mmoriarity/detect-presence-srv:latest"
}

data "docker_registry_image" "go_links_srv" {
  name = "mmoriarity/go-links-srv:latest"
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

data "docker_registry_image" "prometheus_backup" {
  name = "mmoriarity/prometheus-backup:latest"
}

data "docker_registry_image" "tripplite_exporter" {
  name = "mmoriarity/tripplite-exporter:latest"
}

data "docker_registry_image" "unifi_exporter" {
  name = "mmoriarity/unifi_exporter:latest"
}

