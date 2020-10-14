{{ define "pushover.priority" }}{{ if eq .Status "firing" }}{{ if eq .CommonLabels.severity "notice" }}0{{ else }}2{{ end }}{{ else }}0{{ end }}{{ end }}

{{ define "pushover.title" }}{{ .CommonAnnotations.summary }}{{ end }}
