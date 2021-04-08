package linksservice

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

var meter = global.Meter("github.com/mjm/pi-tools/go-links/service/linksservice")

var linksCreatedTotal = metric.Must(meter).NewInt64Counter("golinks.links.created.total")
