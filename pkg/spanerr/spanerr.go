package spanerr

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func RecordError(ctx context.Context, err error) error {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
	return err
}
