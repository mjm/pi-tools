package linksservice

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("github.com/mjm/pi-tools/go-links/service/linksservice")

var linksCreatedTotal = metric.Must(meter).NewInt64Counter("golinks.links.created.total")
