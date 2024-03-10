package coder_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

const (
	testContextKey   = "test"
	testContextValue = "test"
)

func TestRegisterDecoder(t *testing.T) {
	t.Run("decoder", func(t *testing.T) {
		testRegisterDecoder(t, func(ctx context.Context, reader io.Reader, target any) error {
			buf, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			reflect.ValueOf(target).Elem().SetString(string(buf) + ctx.Value(testContextKey).(string))

			return nil
		}, "decoder", "decoder-alias")
	})

	t.Run("decoder-without-context", func(t *testing.T) {
		testRegisterDecoder(t, func(reader io.Reader, target any) error {
			buf, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			reflect.ValueOf(target).Elem().SetString(string(buf) + testContextValue)

			return nil
		}, "decoder-without-context", "decoder-without-context-alias")
	})
}

func testRegisterDecoder[D coder.DecoderConstraint](t *testing.T, decoder D, name, alias string) {
	t.Helper()

	coder.RegisterDecoder(decoder, name, alias)

	for _, key := range []string{name, alias} {
		dec := coder.GetDecoder(key)
		if dec == nil {
			t.Error("decoder not found")
			continue
		}

		reader := strings.NewReader(key)
		ctx := context.Background()
		ctx = context.WithValue(ctx, testContextKey, testContextValue)

		var str string

		if err := dec(ctx, reader, &str); err != nil {
			t.Error(err)
		} else if str != key+testContextValue {
			t.Error("should decode data")
		}
	}
}
