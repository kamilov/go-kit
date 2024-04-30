package text_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/transport/http/coder/text"
)

type (
	marshaler struct {
		str string
		err error
	}

	stringer string
)

func (m marshaler) MarshalText() ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}

	return []byte(m.str), nil
}

func (s stringer) String() string {
	return "_" + string(s) + "_"
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		input  any
		target string
	}{
		{"test", "test"},
		{[]byte("test"), "test"},
		{5, "5"},
		{int8(5), "5"},
		{int16(5), "5"},
		{int32(5), "5"},
		{int64(5), "5"},
		{uint(5), "5"},
		{uint8(5), "5"},
		{uint16(5), "5"},
		{uint32(5), "5"},
		{uint64(5), "5"},
		{float32(5.11), "5.11"},
		{5.11, "5.11"},
		{true, "true"},
		{false, "false"},
		{errors.New("error"), "error"},
		{marshaler{str: "test"}, "test"},
		{&marshaler{str: "test"}, "test"},
		{stringer("test"), "_test_"},
	}

	enc := coder.GetEncoder("text")
	if enc == nil {
		t.Fatal("encoder not found")
	}

	for _, test := range tests {
		ctx := context.Background()
		buf := bytes.NewBufferString("")

		if err := enc(ctx, buf, test.input); err != nil {
			t.Error(err)
		} else {
			if diff := cmp.Diff(test.target, buf.String()); diff != "" {
				t.Errorf("%T: %s", test.input, diff)
			}
		}
	}
}

func TestEncoderError(t *testing.T) {
	tests := []struct {
		input any
		err   error
	}{
		{&struct{}{}, text.ErrUnsupportedType},
		{marshaler{err: io.EOF}, io.EOF},
	}

	enc := coder.GetEncoder("text")
	if enc == nil {
		t.Fatal("encoder not found")
	}

	for _, test := range tests {
		ctx := context.Background()
		buf := bytes.NewBufferString("")

		if err := enc(ctx, buf, test.input); err == nil {
			t.Error("should return error")
		} else if !errors.Is(err, test.err) {
			t.Errorf("error is %v, want %v", err, test.err)
		}
	}
}
