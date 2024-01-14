package env

type (
	optionFunc func(loader *Loader)

	// Option interface for configuration types
	Option interface {
		apply(loader *Loader)
	}
)

func (fn optionFunc) apply(loader *Loader) {
	fn(loader)
}

// WithPrefix configuration func to override prefix on Loader
func WithPrefix(prefix string) Option {
	return optionFunc(func(loader *Loader) {
		loader.prefix = prefix
	})
}

// WithLookup configuration func to override lookup env values on Loader
func WithLookup(lookup LookupFunc) Option {
	return optionFunc(func(loader *Loader) {
		loader.lookup = lookup
	})
}
