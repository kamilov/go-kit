package coder_test

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

type contextKey string

const (
	testContextKey   contextKey = "testContextKey"
	testContextValue string     = "testContextValue"
)

func testRegisterDecoder[T coder.DecoderConstraint](t *testing.T, name, alias string, decoder T) {
	t.Helper()

	coder.RegisterDecoder(decoder, name, alias)

	for _, key := range []string{name, alias} {
		dec := coder.GetDecoder(key)
		if dec == nil {
			t.Errorf("decoder for %s should not be nil", key)
			continue
		}

		reader := strings.NewReader(key)
		ctx := context.Background()
		ctx = context.WithValue(ctx, testContextKey, testContextValue)

		var result string

		if err := dec(ctx, reader, &result); err != nil {
			t.Error(err)
		} else if result != key+testContextValue {
			t.Error("should decode data")
		}
	}
}

func TestRegisterDecoder(t *testing.T) {
	t.Run("simple register decoder", func(t *testing.T) {
		testRegisterDecoder(t, "decoder", "decoder-alias", func(ctx context.Context, reader io.Reader, data any) error {
			buf, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			reflect.ValueOf(data).Elem().SetString(string(buf) + ctx.Value(testContextKey).(string))

			return nil
		})
	})

	t.Run("decoder without context", func(t *testing.T) {
		testRegisterDecoder(t, "decoder", "decoder-alias", func(reader io.Reader, data any) error {
			buf, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			reflect.ValueOf(data).Elem().SetString(string(buf) + testContextValue)

			return nil
		})
	})
}
