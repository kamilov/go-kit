package coder

import (
	"context"
	"io"
)

type (
	Decoder               = func(ctx context.Context, reader io.Reader, data any) error
	DecoderWithoutContext = func(reader io.Reader, target any) error

	DecoderConstraint interface {
		Decoder | DecoderWithoutContext
	}
)

var decoders = map[string]Decoder{}

// RegisterDecoder registers a data decoder for a given type name and type aliases
func RegisterDecoder[D DecoderConstraint](decoder D, name string, aliases ...string) {
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

// GetDecoder return the decoder for a given type name
func GetDecoder(name string) Decoder {
	return decoders[name]
}
