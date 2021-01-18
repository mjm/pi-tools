job "backup-borg" {
  datacenters = ["dc1"]

  type = "batch"

  // Backups should be able to preempt normal service workloads, since they have to run on one node
  // and thus are likely to get stuck otherwise.
  priority = 70

  periodic {
    cron             = "30 */4 * * *"
    prohibit_overlap = true
  }

  group "backup" {
    count = 1

    volume "prometheus_data" {
      type   = "host"
      source = "prometheus_data"
    }

    volume "homelab_nfs" {
      type   = "host"
      source = "homelab_nfs"
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
    }

    task "prometheus-snapshot" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image   = "mmoriarity/prometheus-backup@__PROMETHEUS_BACKUP_DIGEST__"
        command = "/prometheus-backup"
        args    = [
          "-prometheus-url",
          "http://127.0.0.1:9090",
          "-prometheus-data-path",
          "/prometheus",
          "-backup-path",
          "${NOMAD_ALLOC_DIR}/data/prometheus",
        ]

        network_mode = "host"
      }

      resources {
        cpu    = 100
        memory = 100
      }

      volume_mount {
        volume      = "prometheus_data"
        destination = "/prometheus"
      }
    }

    task "dump-presence-db" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image   = "postgres@sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"
        command = "pg_dump"
        args    = [
          "--host=postgresql.service.consul",
          "--dbname=presence",
          "--file=${NOMAD_ALLOC_DIR}/data/presence.sql",
        ]

        network_mode = "host"
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["presence"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/presence" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
      }
    }

    task "dump-golinks-db" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image   = "postgres@sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"
        command = "pg_dump"
        args    = [
          "--host=postgresql.service.consul",
          "--dbname=golinks",
          "--file=${NOMAD_ALLOC_DIR}/data/golinks.sql",
        ]

        network_mode = "host"
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["go-links"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/go-links" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
      }
    }

    task "dump-homebase-bot-db" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image   = "postgres@sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"
        command = "pg_dump"
        args    = [
          "--host=postgresql.service.consul",
          "--dbname=homebase_bot",
          "--file=${NOMAD_ALLOC_DIR}/data/homebase_bot.sql",
        ]

        network_mode = "host"
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["homebase-bot"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/homebase-bot" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
      }
    }

    task "dump-grafana-db" {
      lifecycle {
        hook = "prestart"
      }

      driver = "docker"

      config {
        image   = "postgres@sha256:b6399aef923e0529a4f2a5874e8860d29cef3726ab7079883f3368aaa2a9f29c"
        command = "pg_dump"
        args    = [
          "--host=postgresql.service.consul",
          "--dbname=grafana",
          "--file=${NOMAD_ALLOC_DIR}/data/grafana.sql",
        ]

        network_mode = "host"
      }

      resources {
        cpu    = 50
        memory = 50
      }

      vault {
        policies = ["grafana"]
      }

      template {
        // language=GoTemplate
        data        = <<EOF
{{ with secret "database/creds/grafana" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env         = true
      }
    }

    task "backup" {
      driver = "docker"

      config {
        image    = "mmoriarity/borgbackup@sha256:66a6c3f35b441b35f111ef7817605be03a3a45e9be58d0186fc100a7e65abf65"
        command  = "borg"
        args     = [
          "create",
          "--stats",
          "/dest/backup::backup-${NOMAD_ALLOC_ID}",
          "data",
        ]
        work_dir = "${NOMAD_ALLOC_DIR}"
      }

      env {
        BORG_UNKNOWN_UNENCRYPTED_REPO_ACCESS_IS_OK = "yes"
      }

      resources {
        cpu    = 100
        memory = 100
      }

      volume_mount {
        volume      = "homelab_nfs"
        destination = "/dest"
      }
    }
  }
}
