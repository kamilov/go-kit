package zap

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	options struct {
		output     io.Writer
		encoder    zapcore.Encoder
		zapOptions []zap.Option
	}

	// Option interface for configuration types
	Option interface {
		apply(opts *options)
	}

	optionFunc func(opts *options)
)

func (fn optionFunc) apply(opts *options) {
	fn(opts)
}

// WithOutput configuration func to override output
func WithOutput(output io.Writer) Option {
	return optionFunc(func(opts *options) {
		opts.output = output
	})
}

// WithEncoder configuration func to override encoder
func WithEncoder(encoder zapcore.Encoder) Option {
	return optionFunc(func(opts *options) {
		opts.encoder = encoder
	})
}

// WithCaller configuration func to enable/disable print caller on log message
func WithCaller(enabled bool) Option {
	return optionFunc(func(opts *options) {
		opts.zapOptions = append(opts.zapOptions, zap.WithCaller(enabled))
	})
}

// WithCallerSkip configuration func to set skip index to caller backtrace
func WithCallerSkip(skip int) Option {
	return optionFunc(func(opts *options) {
		opts.zapOptions = append(opts.zapOptions, zap.AddCallerSkip(skip))
	})
}

// WithClock configuration func to override clock
func WithClock(clock zapcore.Clock) Option {
	return optionFunc(func(opts *options) {
		opts.zapOptions = append(opts.zapOptions, zap.WithClock(clock))
	})
}

// WithField defining in the configuration the fields that will be displayed in all messages
func WithField(key, value string) Option {
	return optionFunc(func(opts *options) {
		opts.zapOptions = append(opts.zapOptions, zap.Fields(zap.String(key, value)))
	})
}
