// Package json provides encoding and decoding of JSON data.
package json

import (
	"context"
	"encoding/json"
	"io"

	"github.com/kamilov/go-kit/coder"
)

type contextKey int

const (
	// AllowUnknownFields Allow unknown fields option for json decoder
	AllowUnknownFields contextKey = iota
	// UseNumber use number type for json decoder
	UseNumber
	// EscapeHTML escapes html tags for json encoder
	EscapeHTML
)

func init() {
	name := "json"
	alias := "application/json"

	coder.RegisterDecoder(decode, name, alias)
	coder.RegisterEncoder(encode, name, alias)
}

func decode(ctx context.Context, reader io.Reader, target any) error {
	dec := json.NewDecoder(reader)

	if disallowUnknownFields(ctx) {
		dec.DisallowUnknownFields()
	}

	if useNumber(ctx) {
		dec.UseNumber()
	}

	return dec.Decode(target)
}

func encode(ctx context.Context, writer io.Writer, data any) error {
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
