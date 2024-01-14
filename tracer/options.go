package tracer

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

type (
	options struct {
		host                string
		port                int
		name                string
		maxQueueSize        int
		batchTimeout        time.Duration
		exportTimeout       time.Duration
		maxExportsBatchSize int
		sampler             trace.Sampler
		attributes          []attribute.KeyValue
	}

	optionFunc func(opts *options)

	// Option options apply interface
	Option interface {
		apply(opts *options)
	}
)

const (
	defaultPort               = 6831
	defaultMaxQueueSize       = 2048
	defaultBatchTimeout       = 5000 * time.Millisecond
	defaultExportTimeout      = 30000 * time.Millisecond
	defaultMaxExportBatchSize = 512
)

func (f optionFunc) apply(opts *options) {
	f(opts)
}

func defaultOptions() *options {
	//nolint:exhaustruct // default options
	return &options{
		port:                defaultPort,
		maxQueueSize:        defaultMaxQueueSize,
		batchTimeout:        defaultBatchTimeout,
		exportTimeout:       defaultExportTimeout,
		maxExportsBatchSize: defaultMaxExportBatchSize,
	}
}

// WithHost configure option jaeger hostname
func WithHost(host string) Option {
	return optionFunc(func(opts *options) {
		opts.host = host
	})
}

// WithPort configure option jaeger port
func WithPort(port int) Option {
	return optionFunc(func(opts *options) {
		opts.port = port
	})
}

// WithName configure option instrumentation name
func WithName(name string) Option {
	return optionFunc(func(opts *options) {
		opts.name = name
	})
}

// WithMaxQueueSize configure option maximum queue size
func WithMaxQueueSize(size int) Option {
	return optionFunc(func(opts *options) {
		opts.maxQueueSize = size
	})
}

// WithBatchTimeout configure option batch timeout time
func WithBatchTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *options) {
		opts.batchTimeout = timeout
	})
}

// WithExportTimeout configure option export timeout time
func WithExportTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *options) {
		opts.exportTimeout = timeout
	})
}

// WithMaxExportBatchSize configure option maximum export batch size
func WithMaxExportBatchSize(size int) Option {
	return optionFunc(func(opts *options) {
		opts.maxExportsBatchSize = size
	})
}

// WithSampler configure option sampler
func WithSampler(sampler trace.Sampler) Option {
	return optionFunc(func(opts *options) {
		opts.sampler = sampler
	})
}

// WithAttribute configure option attributes. Add tag to list
func WithAttribute(tag attribute.KeyValue) Option {
	return optionFunc(func(opts *options) {
		opts.attributes = append(opts.attributes, tag)
	})
}
