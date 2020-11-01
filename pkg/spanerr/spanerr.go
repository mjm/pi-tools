package spanerr

import (
	"context"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
)

func RecordError(ctx context.Context, err error) error {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
	return err
}
