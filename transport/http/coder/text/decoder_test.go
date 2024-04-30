package text_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/transport/http/coder/text"
)

type unmarshaler struct {
	Str string
	Err error
}

func (m *unmarshaler) UnmarshalText(data []byte) error {
	if m.Err != nil {
		return m.Err
	}

	m.Str = string(data)

	return nil
}

func pointer[T any](v T) *T {
	return &v
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		input  string
		target any
	}{
		{" \n", ""},
		{"test\n", "test"},
		{"test\n", []byte("test")},
		{"5\n", 5},
		{"5\n", int8(5)},
		{"5\n", int16(5)},
		{"5\n", int32(5)},
		{"5\n", int64(5)},
		{"5\n", uint(5)},
		{"5\n", uint8(5)},
		{"5\n", uint16(5)},
		{"5\n", uint32(5)},
		{"5\n", uint64(5)},
		{"5.12\n", float32(5.12)},
		{"5.12\n", 5.12},
		{"TRUE", true},
		{"FALSE", false},
		{"test\n", unmarshaler{Str: "test"}},
		{"test\n", unmarshaler{Str: "test"}},
		{"test\n", &unmarshaler{Str: "test"}},
	}

	dec := coder.GetDecoder("text")
	if dec == nil {
		t.Error("decoder not found")
	}

	for _, test := range tests {
		ctx := context.Background()
		val := reflect.New(reflect.TypeOf(test.target)).Interface()
		buf := bytes.NewBufferString(test.input)

		if err := dec(ctx, buf, val); err != nil {
			t.Error(err)
		} else {
			if diff := cmp.Diff(test.target, reflect.ValueOf(val).Elem().Interface()); diff != "" {
				t.Errorf("%T, %s", test.target, diff)
			}
		}
	}
}

func TestDecoderError(t *testing.T) {
	tests := []struct {
		input  string
		target any
		err    error
	}{
		{"test\n", &struct{}{}, text.ErrUnsupportedType},
		{"test\n", pointer(0), strconv.ErrSyntax},
		{"test\n", &unmarshaler{Err: io.EOF}, io.EOF},
	}

	dec := coder.GetDecoder("text")
	if dec == nil {
		t.Error("decoder not found")
	}

	for _, test := range tests {
		ctx := context.Background()
		buf := bytes.NewBufferString(test.input)

		if err := dec(ctx, buf, test.target); err == nil {
			t.Error("should return error")
		} else if !errors.Is(err, test.err) {
			t.Errorf("should return error %v, got %v", test.err, err)
		}
	}
}
