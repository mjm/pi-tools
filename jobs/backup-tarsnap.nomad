locals {
  databases = {
    presence       = {
      db_name    = "presence"
      vault_role = "presence"
    }
    golinks        = {
      db_name    = "golinks"
      vault_role = "go-links"
    }
    "homebase-bot" = {
      db_name    = "homebase_bot"
      vault_role = "homebase-bot"
    }
    grafana        = {
      db_name    = "grafana"
      vault_role = "grafana"
    }
    paperless      = {
      db_name    = "paperless"
      vault_role = "paperless"
    }
  }

  phabricator_databases = [
    "phabricator_almanac",
    "phabricator_application",
    "phabricator_audit",
    "phabricator_auth",
    "phabricator_badges",
    "phabricator_cache",
    "phabricator_calendar",
    "phabricator_chatlog",
    "phabricator_conduit",
    "phabricator_config",
    "phabricator_conpherence",
    "phabricator_countdown",
    "phabricator_daemon",
    "phabricator_dashboard",
    "phabricator_differential",
    "phabricator_diviner",
    "phabricator_doorkeeper",
    "phabricator_draft",
    "phabricator_drydock",
    "phabricator_fact",
    "phabricator_feed",
    "phabricator_file",
    "phabricator_flag",
    "phabricator_fund",
    "phabricator_harbormaster",
    "phabricator_herald",
    "phabricator_legalpad",
    "phabricator_maniphest",
    "phabricator_meta_data",
    "phabricator_metamta",
    "phabricator_multimeter",
    "phabricator_nuance",
    "phabricator_oauth_server",
    "phabricator_owners",
    "phabricator_packages",
    "phabricator_passphrase",
    "phabricator_paste",
    "phabricator_pastebin",
    "phabricator_phame",
    "phabricator_phlux",
    "phabricator_pholio",
    "phabricator_phortune",
    "phabricator_phragment",
    "phabricator_phrequent",
    "phabricator_phriction",
    "phabricator_phurl",
    "phabricator_policy",
    "phabricator_ponder",
    "phabricator_project",
    "phabricator_releeph",
    "phabricator_repository",
    "phabricator_search",
    "phabricator_slowvote",
    "phabricator_spaces",
    "phabricator_system",
    "phabricator_token",
    "phabricator_user",
    "phabricator_worker",
    "phabricator_xhpast",
    "phabricator_xhprof",
  ]
}

job "backup-tarsnap" {
  datacenters = ["dc1"]

  type = "batch"

  // Backups should be able to preempt normal service workloads, since they have to run on one node
  // and thus are likely to get stuck otherwise.
  priority = 70

  periodic {
    cron             = "0 12 * * *"
    prohibit_overlap = true
  }

  meta {
    logging_tag = "backup-tarsnap"
  }

  group "backup" {
    count = 1

    constraint {
      attribute = "${node.unique.name}"
      value     = "raspberrypi"
    }

    task "consul-snapshot" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image   = "consul@sha256:7b878010be55876f2dd419e0e95aad54cd87ae078d5de54e232e4135eb1069c6"
        command = "/bin/sh"
        args    = ["-c", "consul snapshot save ${NOMAD_ALLOC_DIR}/data/consul.snap"]

        network_mode = "host"
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["backup"]
      }

      template {
        data        = <<EOF
CONSUL_HTTP_TOKEN={{ with secret "consul/creds/backup" }}{{ .Data.token }}{{ end }}
EOF
        destination = "secrets/consul.env"
        env         = true
      }
    }

    dynamic "task" {
      for_each = local.databases

      labels = ["dump-${task.key}-db"]

      content {
        lifecycle {
          hook = "prestart"
        }

        driver = "docker"

        config {
          image   = "postgres@sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"
          command = "pg_dump"
          args    = [
            "--host=10.0.2.102",
            "--dbname=${DB_NAME}",
            "--file=${NOMAD_ALLOC_DIR}/data/${DB_NAME}.sql",
          ]

          network_mode = "host"
        }

        resources {
          cpu    = 50
          memory = 50
        }

        vault {
          policies = [task.value.vault_role]
        }

        template {
          // language=GoTemplate
          data        = <<EOF
{{ with secret "database/creds/${task.value.vault_role}" }}
DB_NAME=${task.value.db_name}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
          destination = "secrets/db.env"
          env         = true
        }
      }
    }

    task "dump-phabricator-dbs" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image = "mysql/mysql-server@sha256:b33c6e23c8678e29a43ae7cad47cd6bbead6e39c911c5a7b2b6d943cd42b2944"
        args  = concat([
          "mysqldump",
          "--defaults-file=${NOMAD_SECRETS_DIR}/my.cnf",
          "--result-file=${NOMAD_ALLOC_DIR}/data/phabricator.sql",
          "--databases",
        ], local.phabricator_databases)
      }

      resources {
        cpu    = 200
        memory = 500
      }

      vault {
        policies = ["phabricator"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
[client]
host = mysql.service.consul
{{ with secret "database/creds/phabricator" -}}
user = {{ .Data.username }}
password = {{ .Data.password }}
{{- end }}
EOF
        destination = "secrets/my.cnf"
      }
    }

    task "backup" {
      driver = "docker"

      config {
        image    = "mmoriarity/perform-backup"
        command  = "/usr/bin/perform-backup"
        args     = ["-kind", "tarsnap"]

        mount {
          type   = "bind"
          target = "/var/lib/tarsnap/cache"
          source = "/var/lib/tarsnap/cache"
        }
      }

      resources {
        cpu    = 100
        memory = 100
      }

      vault {
        policies = ["tarsnap"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/tarsnap" }}{{ .Data.data.key | base64Decode }}{{ end }}
EOF
        destination = "secrets/tarsnap.key"
      }

      template {
        // language=GoTemplate
        data        = <<EOF
PUSHGATEWAY_URL={{ with service "pushgateway" }}{{ with index . 0 }}http://{{ .Address }}:{{ .Port }}{{ end }}{{ end }}
EOF
        destination = "local/backup.env"
        env         = true
      }
    }

    task "prune" {
      lifecycle {
        hook = "poststop"
      }

      driver = "docker"

      config {
        image   = "mmoriarity/perform-backup"
        command = "${NOMAD_TASK_DIR}/prune.sh"

        mount {
          type   = "bind"
          target = "/var/lib/tarsnap/cache"
          source = "/var/lib/tarsnap/cache"
        }
      }

      resources {
        cpu    = 200
        memory = 30
      }

      vault {
        policies = ["tarsnap"]
      }

      template {
        data        = <<EOF
{{ with secret "kv/tarsnap" }}{{ .Data.data.key | base64Decode }}{{ end }}
EOF
        destination = "secrets/tarsnap.key"
      }

      template {
        data        = file("backup-tarsnap/prune.sh")
        destination = "local/prune.sh"
        perms       = "0755"
      }
    }
  }
}
