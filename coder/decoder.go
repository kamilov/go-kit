package coder

import (
	"context"
	"io"
)

type (
	Decoder               = func(ctx context.Context, reader io.Reader, data any) error
	DecoderWithoutContext = func(reader io.Reader, data any) error

	DecoderConstraint interface {
		Decoder | DecoderWithoutContext
	}
)

//nolint:gochecknoglobals // used to register decoders
var decoders = map[string]Decoder{}

func RegisterDecoder[T DecoderConstraint](decoder T, name string, aliases ...string) {
	switch d := any(decoder).(type) {
	case DecoderWithoutContext:
		decoders[name] = func(_ context.Context, reader io.Reader, data any) error {
			return d(reader, data)
		}

	case Decoder:
		decoders[name] = d
	}

	for _, alias := range aliases {
		decoders[alias] = decoders[name]
	}
}

func GetDecoder(name string) Decoder {
	return decoders[name]
}
