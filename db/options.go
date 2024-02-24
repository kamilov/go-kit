package db

import (
	"net/url"
	"strings"

	logHook "github.com/kamilov/go-kit/db/hook/log"
	tracerHook "github.com/kamilov/go-kit/db/hook/tracer"
	"github.com/kamilov/go-kit/log"
	"github.com/kamilov/go-kit/tracer"
	"github.com/loghole/dbhook"
)

type (
	// Hook internal hook alias
	Hook = dbhook.Hook

	options struct {
		config *Config
		hooks  []Hook
	}
	optionFunc func(opts *options)

	// Option interface for configuration types
	Option interface {
		apply(opts *options)
	}
)

func (fn optionFunc) apply(opts *options) {
	fn(opts)
}

// WithConfig register database configuration to options
func WithConfig(config *Config) Option {
	return optionFunc(func(opts *options) {
		opts.config = config
	})
}

// WithConfigFromDSN convert connection string to configuration struct and register option
func WithConfigFromDSN(dsn string) Option {
	u, _ := url.Parse(dsn)
	//nolint:exhaustruct // define known fields
	config := &Config{
		Hostname: u.Host,
		Username: u.User.Username(),
		Database: u.Path,
		Driver:   driverName(u.Scheme),
	}

	if password, set := u.User.Password(); set {
		config.Password = password
	}

	for key, value := range u.Query() {
		config.Params[key] = strings.Join(value, ",")
	}

	return WithConfig(config)
}

// WithHook add hook to database working process
func WithHook(hook Hook) Option {
	return optionFunc(func(opts *options) {
		opts.hooks = append(opts.hooks, hook)
	})
}

// WithLogHook log queries
func WithLogHook(logger log.Logger) Option {
	return WithHook(logHook.New(logger))
}

// WithTracerHook opentracing hook
func WithTracerHook(t tracer.Tracer) Option {
	return WithHook(tracerHook.New(t))
}
