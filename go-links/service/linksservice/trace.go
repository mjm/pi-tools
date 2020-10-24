package linksservice

import (
	"go.opentelemetry.io/otel/api/global"
)

var tracer = global.Tracer("github.com/mjm/pi-tools/go-links/service/linksservice")
