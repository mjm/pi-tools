package linksservice

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
)

var meter = global.Meter("github.com/mjm/pi-tools/go-links/service/linksservice")

var linksCreatedTotal = metric.Must(meter).NewInt64Counter("golinks.links.created.total")
