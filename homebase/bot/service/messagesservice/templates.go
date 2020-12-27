package messagesservice

import (
	"strings"
	"text/template"
	"time"

	"github.com/hako/durafmt"
	"github.com/jonboulle/clockwork"
)

const (
	tripCompletedTemplate = "TripCompleted"
)

const templatesText = `
{{ define "TripCompleted" -}}
You {{if lt (.ReturnedAt | ago) (duraparse "5m") -}}
just returned
{{- else -}}
returned {{ .ReturnedAt | ago | durafmtshort }} ago
{{- end }} from a trip that lasted *{{ .Duration | durafmtshort }}*\.
{{- if gt (len .Tags) 0 }}

🏷 {{ .Tags | join ", " }}
{{- end }}
{{- end }}
`

type tripCompletedTemplateInput struct {
	ReturnedAt time.Time
	Duration   time.Duration
	Tags       []string
}

var templates *template.Template

func init() {
	templates = template.Must(parseTemplates(clockwork.NewRealClock()))
}

func parseTemplates(clock clockwork.Clock) (*template.Template, error) {
	t := template.New("templates")

	t.Funcs(map[string]interface{}{
		"ago": func(t time.Time) time.Duration {
			return clock.Now().Sub(t)
		},
		"durafmtshort": durafmt.ParseShort,
		"duraparse":    time.ParseDuration,
		"join": func(sep string, elems []string) string {
			return strings.Join(elems, sep)
		},
	})

	return t.Parse(templatesText)
}
