{{ define "pushover.priority" }}{{ if eq .Status "firing" }}{{ if eq .CommonLabels.severity "notice" }}0{{ else }}2{{ end }}{{ else }}0{{ end }}{{ end }}

{{ define "pushover.title" }}[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .CommonAnnotations.summary }}{{ end }}
