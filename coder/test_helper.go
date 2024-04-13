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

type TestOptions[T any] struct {
	Name    string
	Encoded string
	Decoded T
}

//nolint:revive // testing arguments sorting
func Test[T any](t *testing.T, ctx context.Context, options TestOptions[T]) {
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

//nolint:revive // testing arguments sorting
func testDecode[T any](t *testing.T, ctx context.Context, options TestOptions[T]) {
	t.Helper()

	decoder := GetDecoder(options.Name)
	if decoder == nil {
		t.Fatalf("decoder for %s should not be nil", options.Name)
	}

	buf := &bytes.Buffer{}

	buf.WriteString(options.Encoded)

	var decoded T

	if err := decoder(ctx, buf, &decoded); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(options.Decoded, decoded) {
		t.Errorf("decoded value does not match the expected value")
	}
}

//nolint:revive // testing arguments sorting
func testEncode[T any](t *testing.T, ctx context.Context, options TestOptions[T]) {
	t.Helper()

	encoder := GetEncoder(options.Name)
	if encoder == nil {
		t.Fatalf("encoder for %s should not be nil", options.Name)
	}

	buf := &bytes.Buffer{}

	if err := encoder(ctx, buf, &options.Decoded); err != nil {
		t.Error(err)
	} else if diff := cmp.Diff(options.Encoded, strings.TrimRight(buf.String(), "\n"), ignoreUnexported[T]()); diff != "" {
		t.Errorf("encoded value does not match the expected value")
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

//nolint:revive // testing arguments sorting
func Benchmark[T any](b *testing.B, ctx context.Context, options TestOptions[T]) {
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

//nolint:revive // testing arguments sorting
func benchmarkDecode[T any](b *testing.B, ctx context.Context, options TestOptions[T]) {
	b.Helper()

	decoder := GetDecoder(options.Name)
	if decoder == nil {
		b.Fatalf("decoder for %s should not be nil", options.Name)
	}

	buf := strings.NewReader(options.Encoded)

	var decoded T

	for i := 0; i < b.N; i++ {
		_, _ = buf.Seek(0, io.SeekStart)
		_ = decoder(ctx, buf, &decoded)
	}
}

//nolint:revive // testing arguments sorting
func benchmarkEncode[T any](b *testing.B, ctx context.Context, options TestOptions[T]) {
	b.Helper()

	encoder := GetEncoder(options.Name)
	if encoder == nil {
		b.Fatalf("encoder for %s should not be nil", options.Name)
	}

	buf := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = encoder(ctx, buf, options.Decoded)
	}
}
