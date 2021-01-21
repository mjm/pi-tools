terraform {
  backend "consul" {
    scheme = "http"
    access_token = ""
    datacenter = "dc1"
    path = "terraform/state"
  }

  required_providers {
    nomad = {
      source = "hashicorp/nomad"
      version = "1.4.12"
    }
    consul = {
      source = "hashicorp/consul"
      version = "2.11.0"
    }
    local = {
      source = "hashicorp/local"
      version = "2.0.0"
    }
    docker = {
      source = "kreuzwerker/docker"
      version = "2.10.0"
    }
  }
}

provider "nomad" {
}

provider "consul" {
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}
