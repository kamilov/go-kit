// Package tracer tracing for jaeger
package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

// Tracer main struct for tracing
type Tracer struct {
	trace.Tracer
	provider *sdktrace.TracerProvider
}

// New creates new Tracer
func New(ctx context.Context, opts ...Option) (*Tracer, error) {
	cfg := defaultOptions()

	for _, opt := range opts {
		opt.apply(cfg)
	}

	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(
			fmt.Sprintf("%s:%d", cfg.host, cfg.port),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("init jaeger exporter: %w", err)
	}

	spanProcessor := sdktrace.NewBatchSpanProcessor(
		exporter,
		sdktrace.WithMaxQueueSize(cfg.maxQueueSize),
		sdktrace.WithBatchTimeout(cfg.batchTimeout),
		sdktrace.WithExportTimeout(cfg.exportTimeout),
		sdktrace.WithMaxExportBatchSize(cfg.maxExportsBatchSize),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resource.NewSchemaless(cfg.attributes...)),
		sdktrace.WithSampler(sdktrace.ParentBased(cfg.sampler)),
	)

	tracer := &Tracer{
		Tracer:   tracerProvider.Tracer(cfg.name),
		provider: tracerProvider,
	}

	return tracer, nil
}

// NewNoop creates nullable Tracer
func NewNoop(name string) *Tracer {
	return &Tracer{
		Tracer:   noop.NewTracerProvider().Tracer(name),
		provider: sdktrace.NewTracerProvider(),
	}
}

// Sync synchronization tracing
func (t *Tracer) Sync() {
	_ = t.provider.Shutdown(context.Background())
}
