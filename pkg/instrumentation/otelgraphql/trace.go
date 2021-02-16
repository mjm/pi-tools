package otelgraphql

import (
	"context"

	"github.com/mjm/graphql-go/errors"
	"github.com/mjm/graphql-go/introspection"
	gqltrace "github.com/mjm/graphql-go/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/mjm/pi-tools/pkg/instrumentation/otelgraphql")

var (
	queryKey         = label.Key("graphql.query")
	operationNameKey = label.Key("graphql.operation_name")
	typeKey          = label.Key("graphql.type")
	fieldKey         = label.Key("graphql.field")
)

// GraphQLTracer implements gqltrace.Tracer for tracing GraphQL requests
type GraphQLTracer struct{}

var _ gqltrace.Tracer = GraphQLTracer{}

func (GraphQLTracer) TraceQuery(ctx context.Context, queryString string, operationName string, variables map[string]interface{}, varTypes map[string]*introspection.Type) (context.Context, gqltrace.TraceQueryFinishFunc) {
	if operationName == "IntrospectionQuery" {
		return ctx, func([]*errors.QueryError) {}
	}

	trace.SpanFromContext(ctx).SetAttributes(operationNameKey.String(operationName))

	ctx, span := tracer.Start(ctx, "graphql.Query",
		trace.WithAttributes(
			queryKey.String(queryString),
			operationNameKey.String(operationName)))

	return ctx, func(errs []*errors.QueryError) {
		for _, err := range errs {
			recordQueryError(span, err)
		}
		span.End()
	}
}

func (GraphQLTracer) TraceField(ctx context.Context, label, typeName, fieldName string, trivial bool, args map[string]interface{}) (context.Context, gqltrace.TraceFieldFinishFunc) {
	if trivial {
		return ctx, func(*errors.QueryError) {}
	}

	if fieldName == "__schema" || typeName == "__Schema" || typeName == "__Type" {
		return ctx, func(*errors.QueryError) {}
	}

	ctx, span := tracer.Start(ctx, label,
		trace.WithAttributes(
			typeKey.String(typeName),
			fieldKey.String(fieldName)))

	return ctx, func(err *errors.QueryError) {
		if err != nil {
			recordQueryError(span, err)
		}
		span.End()
	}
}

func recordQueryError(span trace.Span, err *errors.QueryError) {
	if err.ResolverError == nil {
		span.RecordError(err)
		return
	}

	span.RecordError(err.ResolverError)
}
