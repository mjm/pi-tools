{{ define "pagerduty.severity" }}{{ if eq .Status "firing" }}{{ if eq .CommonLabels.severity "notice" }}warning{{ else }}error{{ end }}{{ else }}warning{{ end }}{{ end }}
