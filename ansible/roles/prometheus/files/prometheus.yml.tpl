global:
  scrape_interval: 60s
  evaluation_interval: 30s

alerting:
  alertmanagers:
    - consul_sd_configs:
        - services: [alertmanager]
          server: 10.0.2.10:8500
          token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}

rule_files:
  - /usr/local/etc/rules/*.yml

scrape_configs:
  - job_name: consul-agent
    consul_sd_configs:
      - services: [consul]
        server: 10.0.2.10:8500
        token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
    bearer_token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
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
        server: 10.0.2.10:8500
        token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
    metrics_path: /v1/metrics
    params:
      format: [prometheus]
    scheme: https
    tls_config:
      ca_file: /usr/local/etc/nomad.ca.crt
      cert_file: /usr/local/etc/nomad.crt
      key_file: /usr/local/etc/nomad.key
    relabel_configs:
      - source_labels: [__meta_consul_node]
        target_label: node_name

  - job_name: vault
    consul_sd_configs:
      - services: [vault]
        tags: [active]  # metrics are only available from the active node
        server: 10.0.2.10:8500
        token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
    metrics_path: /v1/sys/metrics
    params:
      format: [prometheus]
    bearer_token: {{ with secret "auth/token/lookup-self" }}{{ .Data.id }}{{ end }}
    relabel_configs:
      - source_labels: [__meta_consul_node]
        target_label: node_name

  - job_name: consul-services
    consul_sd_configs:
      - token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
        server: 10.0.2.10:8500
    relabel_configs:
      - source_labels: [__meta_consul_service]
        action: drop
        regex: (.+)-sidecar-proxy
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
      - source_labels: [__address__, __meta_consul_service_metadata_metrics_port]
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: ${1}:${2}
        target_label: __address__

  - job_name: consul-connect-envoy
    scrape_interval: 300s
    consul_sd_configs:
      - token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
        server: 10.0.2.10:8500
    relabel_configs:
      - source_labels: [__meta_consul_service]
        action: drop
        regex: (.+)-sidecar-proxy
      - source_labels: [__meta_consul_service_metadata_envoy_metrics_port]
        action: keep
        regex: (.+)
      - source_labels: [__meta_consul_service]
        target_label: service_name
        regex: (.+)
      - source_labels: [__meta_consul_node]
        target_label: node_name
        regex: (.+)
      - source_labels: [__address__, __meta_consul_service_metadata_envoy_metrics_port]
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: ${1}:${2}
        target_label: __address__

  - job_name: pushgateway
    consul_sd_configs:
      - services: [pushgateway]
        server: 10.0.2.10:8500
        token: {{ with secret "consul/creds/prometheus" }}{{ .Data.token }}{{ end }}
    honor_labels: true

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
          - 10.0.2.101
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
          - 10.0.2.101
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
          - 10.0.2.101
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

  - job_name: 'blackbox-dns-ad-blocking'
    metrics_path: /probe
    params:
      module: [dns_ad_blocking]
    static_configs:
      - targets:
          - 10.0.2.101
        labels:
          probe_type: dns
          probe_scope: ad-blocking
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: {{ with service "blackbox-exporter" }}{{ with index . 0 }}{{ .Address }}:{{ .Port }}{{ end }}{{ end }}
