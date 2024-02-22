package tracer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/kamilov/go-kit/tracer"
	"github.com/loghole/dbhook"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Hook struct {
	tracer tracer.Tracer
}

func New(tracer tracer.Tracer) *Hook {
	return &Hook{tracer}
}

func (h *Hook) Before(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	parent := trace.SpanFromContext(ctx)
	if parent == nil {
		return ctx, nil
	}

	ctx, span := h.tracer.Start(ctx, spew.Sprintf("SQL[%s]", input.Caller), trace.WithSpanKind(trace.SpanKindInternal))

	span.SetAttributes(
		semconv.DBStatementKey.String(input.Query),
	)

	return trace.ContextWithSpan(ctx, span), nil
}

func (h *Hook) After(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	if span := trace.SpanFromContext(ctx); span != nil {
		defer span.End()
	}
	return ctx, nil
}

func (h *Hook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	if span := trace.SpanFromContext(ctx); span != nil {
		defer span.End()

		if ctx.Err() != nil && errors.Is(ctx.Err(), context.Canceled) {
			return ctx, input.Error
		}

		if input.Error != nil || errors.Is(input.Error, sql.ErrNoRows) {
			return ctx, input.Error
		}

		span.RecordError(input.Error)
		span.SetStatus(codes.Error, "error")
	}
	return ctx, input.Error
}
