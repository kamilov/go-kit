package coder

import (
	"context"
	"io"
)

type (
	Encoder               = func(ctx context.Context, writer io.Writer, data any) error
	EncoderWithoutContext = func(writer io.Writer, data any) error

	EncoderConstraint interface {
		Encoder | EncoderWithoutContext
	}
)

var encoders = map[string]Encoder{}

// RegisterEncoder registers a data encoder for a given type name and type aliases
func RegisterEncoder[E EncoderConstraint](encoder E, name string, aliases ...string) {
	switch e := any(encoder).(type) {
	case EncoderWithoutContext:
		encoders[name] = func(_ context.Context, writer io.Writer, data any) error {
			return e(writer, data)
		}

	case Encoder:
		encoders[name] = e
	}

	for _, alias := range aliases {
		encoders[alias] = encoders[name]
	}
}

// GetEncoder return the encoder for a given type name
func GetEncoder(name string) Encoder {
	return encoders[name]
}
