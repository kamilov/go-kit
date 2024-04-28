package db

import (
	"net/url"
	"strings"

	"github.com/loghole/dbhook"
)

type (
	options struct {
		config      *Config
		hookOptions []dbhook.HookOption
	}

	optionFunc func(*options)

	Option interface {
		apply(*options)
	}
)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithConfig(config *Config) Option {
	return optionFunc(func(o *options) {
		o.config = config
	})
}

func WithConfigDSN(dsn string) Option {
	u, _ := url.Parse(dsn)
	config := &Config{
		Hostname: u.Host,
		Username: u.User.Username(),
		Database: u.Path,
		Params:   make(map[string]string),
	}

	switch u.Scheme {
	case "sqlite":
		config.Driver = SQLite
	default:
		config.Driver = driverName(u.Scheme)
	}

	if password, set := u.User.Password(); set {
		config.Password = password
	}

	for key, value := range u.Query() {
		config.Params[key] = strings.Join(value, ",")
	}

	return WithConfig(config)
}

func WithHook(hooks ...dbhook.Hook) Option {
	return optionFunc(func(o *options) {
		o.hookOptions = append(o.hookOptions, dbhook.WithHook(hooks...))
	})
}

func WithHookBefore(hooks ...dbhook.HookBefore) Option {
	return optionFunc(func(o *options) {
		o.hookOptions = append(o.hookOptions, dbhook.WithHooksBefore(hooks...))
	})
}

func WithHookAfter(hooks ...dbhook.HookAfter) Option {
	return optionFunc(func(o *options) {
		o.hookOptions = append(o.hookOptions, dbhook.WithHooksAfter(hooks...))
	})
}

func WithHookError(hooks ...dbhook.HookError) Option {
	return optionFunc(func(o *options) {
		o.hookOptions = append(o.hookOptions, dbhook.WithHooksError(hooks...))
	})
}
