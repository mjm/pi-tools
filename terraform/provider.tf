terraform {
  backend "consul" {
    scheme = "http"
    access_token = ""
    datacenter = "dc1"
    path = "terraform/state"
  }

  required_providers {
    consul = {
      source = "hashicorp/consul"
      version = "2.11.0"
    }
    vault = {
      source = "hashicorp/vault"
      version = "2.18.0"
    }
  }
}

provider "consul" {
}

provider "vault" {
}
