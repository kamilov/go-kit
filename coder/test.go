// Package coder serialize and unserializer data
package coder

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Options test options
type Options[T any] struct {
	Name    string
	Encoded string
	Decoded T
}

// Test helper a testing encoders and decoder for a given types
func Test[T any](t *testing.T, ctx context.Context, options Options[T]) {
	t.Helper()

	t.Run("decode", func(t *testing.T) {
		t.Helper()
		testDecode(t, ctx, options)
	})

	t.Run("encode", func(t *testing.T) {
		t.Helper()
		testEncode(t, ctx, options)
	})
}

func testDecode[T any](t *testing.T, ctx context.Context, options Options[T]) {
	t.Helper()

	decoder := GetDecoder(options.Name)
	if decoder == nil {
		t.Fatal("decoder not found")
	}

	buf := &bytes.Buffer{}

	buf.WriteString(options.Encoded)

	var decoded T

	if err := decoder(ctx, buf, &decoded); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(options.Decoded, decoded) {
		t.Errorf("should be decode data: %v != %v", options.Decoded, decoded)
	}
}

func testEncode[T any](t *testing.T, ctx context.Context, options Options[T]) {
	t.Helper()

	encoder := GetEncoder(options.Name)
	if encoder == nil {
		t.Fatal("encoder not found")
	}

	buf := &bytes.Buffer{}

	if err := encoder(ctx, buf, options.Decoded); err != nil {
		t.Error(err)
	} else if diff := cmp.Diff(options.Encoded, strings.TrimRight(buf.String(), "\n"), ignoreUnexported[T]()); diff != "" {
		t.Errorf("should be encode data: %s", diff)
	}
}

func ignoreUnexported[T any]() cmp.Option {
	rt := reflect.TypeOf(*new(T))

	for rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		return nil
	}

	return cmpopts.IgnoreUnexported(reflect.New(rt).Elem().Interface())
}

// Benchmark helper a benchmarking encoders and decoder for a given types
func Benchmark[T any](b *testing.B, ctx context.Context, options Options[T]) {
	b.Helper()

	b.Run("decode", func(b *testing.B) {
		b.Helper()
		benchmarkDecode(b, ctx, options)
	})

	b.Run("encode", func(b *testing.B) {
		b.Helper()
		benchmarkEncode(b, ctx, options)
	})
}

func benchmarkDecode[T any](b *testing.B, ctx context.Context, options Options[T]) {
	b.Helper()

	decoder := GetDecoder(options.Name)
	if decoder == nil {
		b.Fatal("decoder not found")
	}

	buf := strings.NewReader(options.Encoded)

	var decoded T

	for i := 0; i < b.N; i++ {
		_, _ = buf.Seek(0, io.SeekStart)
		_ = decoder(ctx, buf, decoded)
	}
}

func benchmarkEncode[T any](b *testing.B, ctx context.Context, options Options[T]) {
	b.Helper()

	encoder := GetEncoder(options.Name)
	if encoder == nil {
		b.Fatal("encoder not found")
	}

	buf := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = encoder(ctx, buf, options.Decoded)
	}
}
