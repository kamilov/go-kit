package coder_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

func testRegisterEncoder[T coder.EncoderConstraint](t *testing.T, name, alias string, encoder T) {
	t.Helper()

	coder.RegisterEncoder(encoder, name, alias)

	for _, key := range []string{name, alias} {
		enc := coder.GetEncoder(key)
		if enc == nil {
			t.Errorf("encoder for %s should not be nil", key)
			continue
		}

		writer := &strings.Builder{}
		ctx := context.Background()
		ctx = context.WithValue(ctx, testContextKey, testContextValue)

		if err := enc(ctx, writer, key); err != nil {
			t.Error(err)
		} else if writer.String() != key+testContextValue {
			t.Error("should encode data")
		}
	}
}

func TestRegisterEncoder(t *testing.T) {
	t.Run("simple encoder", func(t *testing.T) {
		testRegisterEncoder(t, "encoder", "encoder-alias", func(ctx context.Context, writer io.Writer, data any) error {
			_, err := io.WriteString(writer, data.(string)+ctx.Value(testContextKey).(string))

			return err
		})
	})

	t.Run("encoder without context", func(t *testing.T) {
		testRegisterEncoder(
			t,
			"encoder-without-context",
			"encoder-without-context-alias",
			func(writer io.Writer, data any) error {
				_, err := io.WriteString(writer, data.(string)+testContextValue)

				return err
			})
	})
}
