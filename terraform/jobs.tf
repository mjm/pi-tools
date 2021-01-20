resource "nomad_job" "named" {
  jobspec = file("${path.module}/jobs/named.nomad")
}

resource "nomad_job" "pihole" {
  jobspec = file("${path.module}/jobs/pihole.nomad")
}

resource "nomad_job" "prometheus" {
  jobspec = file("${path.module}/jobs/prometheus.nomad")
}

resource "nomad_job" "loki" {
  jobspec = file("${path.module}/jobs/loki.nomad")
}

resource "nomad_job" "promtail" {
  jobspec = file("${path.module}/jobs/promtail.nomad")
}

resource "nomad_job" "node_exporter" {
  jobspec = file("${path.module}/jobs/node-exporter.nomad")
}

resource "nomad_job" "blackbox_exporter" {
  jobspec = file("${path.module}/jobs/blackbox-exporter.nomad")
}

resource "nomad_job" "tripplite_exporter" {
  jobspec = templatefile("${path.module}/jobs/tripplite-exporter.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "unifi_exporter" {
  jobspec = templatefile("${path.module}/jobs/unifi-exporter.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "jaeger" {
  jobspec = file("${path.module}/jobs/jaeger.nomad")
}

resource "nomad_job" "postgresql" {
  jobspec = file("${path.module}/jobs/postgresql.nomad")
}

resource "nomad_job" "grafana" {
  jobspec = file("${path.module}/jobs/grafana.nomad")
}

resource "nomad_job" "deploy" {
  jobspec = templatefile("${path.module}/jobs/deploy.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "beacon_srv" {
  jobspec = templatefile("${path.module}/jobs/beacon-srv.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "presence" {
  jobspec = templatefile("${path.module}/jobs/presence.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "go_links_srv" {
  jobspec = templatefile("${path.module}/jobs/go-links-srv.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "homebase" {
  jobspec = templatefile("${path.module}/jobs/homebase.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "oauth_proxy" {
  jobspec = file("${path.module}/jobs/oauth-proxy.nomad")
}

resource "nomad_job" "ingress" {
  jobspec = file("${path.module}/jobs/ingress.nomad")
  hcl2 {
    enabled = true
  }
}

resource "nomad_job" "backup_tarsnap" {
  jobspec = templatefile("${path.module}/jobs/backup-tarsnap.nomad", {
    image_digests = local.image_digests
  })
}

resource "nomad_job" "backup_borg" {
  jobspec = templatefile("${path.module}/jobs/backup-borg.nomad", {
    image_digests = local.image_digests
  })
}
