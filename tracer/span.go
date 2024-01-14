package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// StartFromContext creates a span and a context.Context use trace provider from parent span
func StartFromContext(
	ctx context.Context,
	tracerName, name string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return trace.SpanFromContext(ctx).
		TracerProvider().
		Tracer(tracerName).
		Start(ctx, name, opts...)
}
