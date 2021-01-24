locals {
  config_files = fileset(".", "named/*")
}

job "named" {
  datacenters = ["dc1"]

  type     = "system"
  priority = 90

  update {
    max_parallel = 1
    stagger      = "30s"
  }

  group "named" {
    network {
      mode = "host"
      port "dns" {
        static = 53
        to     = 53
      }
    }

    service {
      name = "named"
      port = "dns"
    }

    task "named" {
      driver = "docker"
      config {
        image        = "eafxx/bind@sha256:9c15e971a7a358a4ba248e02154b7d5a6b37803bdf65371325364f3cbae9dd43"
        ports        = ["dns"]
        network_mode = "host"
      }

      env {
        WEBMIN_ENABLED = "false"
        DATA_DIR       = "${NOMAD_TASK_DIR}"
      }

      resources {
        cpu    = 50
        memory = 200
      }

      dynamic "template" {
        for_each = local.config_files

        content {
          data          = file(template.value)
          destination   = "local/bind/etc/${basename(template.value)}"
          change_mode   = "signal"
          change_signal = "SIGHUP"
        }
      }
    }
  }
}
