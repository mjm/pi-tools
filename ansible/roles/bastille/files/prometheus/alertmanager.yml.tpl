global:
  resolve_timeout: 5m

templates:
  - /usr/local/etc/alertmanager/templates/*.tpl

route:
  group_by: ['alertname', 'severity']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: pagerduty

receivers:
  - name: pagerduty
    pagerduty_configs:
      - routing_key: << with secret "kv/pagerduty" >><< .Data.data.routing_key >><< end >>
        severity: '{{ template "pagerduty.severity" . }}'

inhibit_rules: []
