package json

import (
	"context"
	"encoding/json"
	"io"

	"github.com/kamilov/go-kit/coder"
)

type contextKey int

const (
	AllowUnknownFields contextKey = iota
	UseNumber
	EscapeHTML
)

//nolint:gochecknoinits // used for automatic adding encode and decode functions
func init() {
	name := "json"
	aliases := []string{"application/json"}

	coder.RegisterDecoder(decoder, name, aliases...)
	coder.RegisterEncoder(encoder, name, aliases...)
}

func decoder(ctx context.Context, reader io.Reader, data any) error {
	dec := json.NewDecoder(reader)

	if disallowUnknownFields(ctx) {
		dec.DisallowUnknownFields()
	}

	if useNumber(ctx) {
		dec.UseNumber()
	}

	return dec.Decode(data)
}

func encoder(ctx context.Context, writer io.Writer, data any) error {
	enc := json.NewEncoder(writer)

	enc.SetEscapeHTML(escapeHTML(ctx))

	return enc.Encode(data)
}

func disallowUnknownFields(ctx context.Context) bool {
	val := ctx.Value(AllowUnknownFields)
	if val == nil {
		return true
	}

	return !val.(bool)
}

func useNumber(ctx context.Context) bool {
	val := ctx.Value(UseNumber)
	if val == nil {
		return false
	}

	return val.(bool)
}

func escapeHTML(ctx context.Context) bool {
	val := ctx.Value(EscapeHTML)
	if val == nil {
		return false
	}

	return val.(bool)
}
