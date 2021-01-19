job "prometheus" {
  datacenters = ["dc1"]

  type     = "service"
  priority = 80

  group "prometheus" {
    count = 1

    volume "data" {
      type      = "host"
      read_only = false
      source    = "prometheus_data"
    }

    network {
      port "http" {
        static = 9090
        to     = 9090
      }
    }

    service {
      name = "prometheus"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }

      check {
        type     = "http"
        path     = "/-/ready"
        interval = "15s"
        timeout  = "3s"
      }
    }

    task "prometheus" {
      driver = "docker"

      config {
        image        = "prom/prometheus@sha256:9fa25ec244e0109fdbeaff89496ac149c0539489f2f2126b9e38cf9837235be4"
        args         = [
          "--web.listen-address=:9090",
          "--config.file=${NOMAD_TASK_DIR}/prometheus.yml",
          "--storage.tsdb.path=/prometheus",
          "--web.console.libraries=/usr/share/prometheus/console_libraries",
          "--web.console.templates=/usr/share/prometheus/consoles",
          "--web.external-url=https://prometheus.homelab/",
          "--web.enable-admin-api",
        ]
        network_mode = "host"

        logging {
          type = "journald"
          config {
            tag = "prometheus"
          }
        }
      }

      resources {
        cpu    = 200
        memory = 1500
      }

      volume_mount {
        volume      = "data"
        destination = "/prometheus"
        read_only   = false
      }

      vault {
        policies      = ["prometheus"]
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        // language=YAML
        data          = <<EOF
global:
  scrape_interval: 30s
  evaluation_interval: 30s

alerting:
  alertmanagers:
    - consul_sd_configs:
        - services: [alertmanager]

rule_files:
  - {{ env "NOMAD_TASK_DIR" }}/rules/*.yml

scrape_configs:
  - job_name: consul-agent
    consul_sd_configs:
      - services: [consul]
    metrics_path: /v1/agent/metrics
    params:
      format: [prometheus]
    relabel_configs:
      - source_labels: [__meta_consul_address]
        target_label: __address__
        replacement: $1:8500
      - source_labels: [__meta_consul_node]
        target_label: node_name

  - job_name: nomad-agent
    consul_sd_configs:
      - services: [nomad-client]
    metrics_path: /v1/metrics
    params:
      format: [prometheus]
    scheme: https
    tls_config:
      ca_file: {{ env "NOMAD_SECRETS_DIR" }}/nomad.ca.crt
      cert_file: {{ env "NOMAD_SECRETS_DIR" }}/nomad.crt
      key_file: {{ env "NOMAD_SECRETS_DIR" }}/nomad.key
    relabel_configs:
      - source_labels: [__meta_consul_node]
        target_label: node_name

  - job_name: vault
    consul_sd_configs:
      - services: [vault]
        tags: [active]  # metrics are only available from the active node
    metrics_path: /v1/sys/metrics
    params:
      format: [prometheus]
    bearer_token_file: {{ env "NOMAD_SECRETS_DIR" }}/vault_token
    relabel_configs:
      - source_labels: [__meta_consul_node]
        target_label: node_name

  - job_name: consul-services
    consul_sd_configs:
      - {}
    relabel_configs:
      - source_labels: [__meta_consul_service_metadata_metrics_path]
        action: keep
        regex: (.+)
      - source_labels: [__meta_consul_service_metadata_metrics_path]
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__meta_consul_service]
        target_label: service_name
        regex: (.+)
      - source_labels: [service_name]
        target_label: service_name
        regex: (.+)-metrics
      - source_labels: [__meta_consul_node]
        target_label: node_name
        regex: (.+)

  - job_name: 'blackbox-dns-public'
    metrics_path: /probe
    params:
      module: [dns_public]
    static_configs:
      - targets:
          - 8.8.4.4
          - 8.8.8.8
          - 1.0.0.1
          - 1.1.1.1
          - 10.0.0.2
          - 10.0.0.3
          - 10.0.0.4
        labels:
          probe_type: dns
          probe_scope: public
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: {{ with service "blackbox-exporter" }}{{ with index . 0 }}{{ .Address }}:{{ .Port }}{{ end }}{{ end }}

  - job_name: 'blackbox-dns-private'
    metrics_path: /probe
    params:
      module: [dns_private]
    static_configs:
      - targets:
          - 10.0.0.2
          - 10.0.0.3
          - 10.0.0.4
        labels:
          probe_type: dns
          probe_scope: private
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: {{ with service "blackbox-exporter" }}{{ with index . 0 }}{{ .Address }}:{{ .Port }}{{ end }}{{ end }}

  - job_name: 'blackbox-dns-private-cname'
    metrics_path: /probe
    params:
      module: [dns_private_cname]
    static_configs:
      - targets:
          - 10.0.0.2
          - 10.0.0.3
          - 10.0.0.4
        labels:
          probe_type: dns
          probe_scope: private-cname
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: {{ with service "blackbox-exporter" }}{{ with index . 0 }}{{ .Address }}:{{ .Port }}{{ end }}{{ end }}
EOF
        destination   = "local/prometheus.yml"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        // language=YAML
        data        = <<EOF
groups:

  - name: dns_alerts
    rules:
      - alert: ExternalDNSFailing
        expr: probe_success{probe_type="dns",instance!~"10.0.0.[2-4]"} < 0.5
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: External DNS queries are failing
          description: >
            One or more external DNS servers are unable to look up Google's home page.
            This might mean your internet is down, though I'm not sure how you're getting this alert
            in that case.

      - alert: DNSFailingExternalSite
        expr: (min(probe_success{probe_type="dns",instance!~"10.0.0.[2-4]"}) > 0.5) and (min(probe_success{probe_type="dns",probe_scope="public",instance=~"10.0.0.[2-4]"}) < 0.5)
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: DNS queries for external sites are failing
          description: >
            The PiHole is unable to look up Google's home page, but external DNS is able to look it up
            without issue. Something is probably wrong with PiHole's configuration.

      - alert: DNSFailingInternalNode
        expr: probe_success{probe_type="dns",probe_scope="private"} < 0.5
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: DNS queries for Raspberry Pi nodes are failing
          description: >
            A DNS server is unable to look up the address for the raspberrypi node, or it's not giving
            the expected IP address.

      - alert: DNSFailingInternalSite
        expr: probe_success{probe_type="dns",probe_scope="private-cname"} < 0.5
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: DNS queries for internal sites are failing
          description: >
            The PiHole is unable to look up the address for Homebase, or it's not returning a CNAME to a
            Raspberry Pi node.

  - name: homebase_bot_alerts
    rules:
      - alert: HomebaseBotServiceDown
        expr: sum(up{app="homebase-bot"}) < 1
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: No homebase-bot-srv instances are running
          description: >
            There don't seem to be any running pods of homebase-bot-srv. The Telegram bot won't work until
            this is addressed.

      - alert: HomebaseBotNoLeader
        expr: sum(homebase_bot_is_leader) < 1
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: homebase-bot-srv has no leader
          description: >
            There are no homebase-bot-srv pods acting as leader. The Telegram bot won't be able to respond to
            incoming messages until there is a leader.

  - name: node_alerts
    rules:
      - alert: NodeExporterDown
        expr: up{app="node-exporter"} < 0.5
        for: 5m
        labels:
          severity: notice
        annotations:
          summary: Node exporter is down
          annotations: >
            One or more of the node-exporters can't be scraped by Prometheus. You should look into this, as you're
            probably missing some metrics you'd like to have.

      - alert: NodeTemperatureTooHigh
        expr: node_hwmon_temp_celsius > 80
        for: 5m
        labels:
          severity: notice
        annotations:
          summary: Raspberry Pi temperature is close to throttling
          description: >
            One or more Raspberry Pis are approaching a temperature that will cause the CPU to start
            throttling.

      - alert: LowDiskSpace
        expr: (node_filesystem_avail_bytes{device!~'rootfs',mountpoint="/"} / node_filesystem_size_bytes{device!~'rootfs',mountpoint="/"}) < 0.05
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: Raspberry Pi is running out of disk space
          description: >
            One or more Raspberry Pis have less than 5% space available on their root volume.

      - alert: HighMemoryUsage
        expr: (node_memory_MemFree_bytes + node_memory_Cached_bytes + node_memory_Buffers_bytes) < 300000000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: Raspberry Pi is using too much memory
          description: >
            One or more Raspberry Pis is using a lot of memory. Some process on the machine probably needs to
            be restarted.

  - name: presence_alerts
    rules:
      - alert: PresenceServiceDown
        expr: sum(up{app="detect-presence"}) < 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: No detect-presence-srv instances are running
          description: >
            There don't seem to be any running pods of detect-presence-srv. Trips won't be able to recorded until this
            is addressed.
EOF
        destination = "local/rules/alerts.yml"
      }

      template {
        // language=YAML
        data        = <<EOF
groups:
  - name: presence
    rules:
      - record: presence:trip_duration_seconds:last
        expr: (presence_last_return_timestamp - presence_last_leave_timestamp) > 0

      - record: presence:trip_duration_seconds:current
        expr: (presence_last_return_timestamp < bool presence_last_leave_timestamp) * (time() - presence_last_leave_timestamp) > 0
EOF
        destination = "local/rules/presence.yml"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.certificate }}
{{ end }}
EOF
        destination   = "secrets/nomad.crt"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.private_key }}
{{ end }}
EOF
        destination   = "secrets/nomad.key"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }

      template {
        data          = <<EOF
{{ with secret "pki-int/issue/nomad-cluster" "ttl=24h" -}}
{{ .Data.issuing_ca }}
{{ end }}
EOF
        destination   = "secrets/nomad.ca.crt"
        change_mode   = "signal"
        change_signal = "SIGHUP"
      }
    }
  }

  group "alertmanager" {
    count = 1

    volume "data" {
      type      = "host"
      read_only = false
      source    = "alertmanager_data"
    }

    network {
      port "http" {
        to = 9093
      }
    }

    service {
      name = "alertmanager"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }

      check {
        type     = "http"
        path     = "/-/ready"
        interval = "15s"
        timeout  = "3s"
      }
    }

    task "alertmanager" {
      driver = "docker"

      config {
        image = "prom/alertmanager@sha256:e690a0f96fcf69c2e1161736d6bb076e22a84841e1ec8ecc87e801c70b942200"
        args  = [
          "--config.file=${NOMAD_SECRETS_DIR}/alertmanager.yml",
          "--storage.path=/alertmanager",
          "--web.external-url=https://alertmanager.homelab/",
        ]
        ports = ["http"]

        logging {
          type = "journald"
          config {
            tag = "alertmanager"
          }
        }
      }

      resources {
        cpu    = 50
        memory = 100
      }

      volume_mount {
        volume      = "data"
        destination = "/alertmanager"
        read_only   = false
      }

      vault {
        policies = [
          "alertmanager",
        ]
      }

      template {
        // language=YAML
        data            = <<EOF
global:
  resolve_timeout: 5m

templates:
  - << env "NOMAD_TASK_DIR" >>/templates/*.tpl

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'push.phone.matt'

receivers:
  - name: 'push.phone.matt'
    pushover_configs:
<<- with secret "kv/alertmanager/pushover" >>
      - user_key: << .Data.data.user_key >>
        token: << .Data.data.token >>
<<- end >>
        title: '{{ template "pushover.title" . }}'
        priority: '{{ template "pushover.priority" . }}'

inhibit_rules: []
EOF
        destination     = "secrets/alertmanager.yml"
        left_delimiter  = "<<"
        right_delimiter = ">>"
      }

      template {
        // language=GoTemplate
        data            = <<EOF
{{ define "pushover.priority" }}{{ if eq .Status "firing" }}{{ if eq .CommonLabels.severity "notice" }}0{{ else }}2{{ end }}{{ else }}0{{ end }}{{ end }}

{{ define "pushover.title" }}[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonAnnotations.summary }}{{ end }}
EOF
        destination     = "local/templates/pushover.tpl"
        left_delimiter  = "<<"
        right_delimiter = ">>"
      }
    }
  }
}
