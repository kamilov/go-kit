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

//nolint:gochecknoglobals // used to register encoders
var encoders = map[string]Encoder{}

func RegisterEncoder[T EncoderConstraint](encoder T, name string, aliases ...string) {
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

func GetEncoder(name string) Encoder {
	return encoders[name]
}
