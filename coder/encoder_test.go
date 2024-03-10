package coder_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

func TestRegisterEncoder(t *testing.T) {
	t.Run("encoder", func(t *testing.T) {
		testRegisterEncoder(t, func(ctx context.Context, writer io.Writer, data any) error {
			_, err := io.WriteString(writer, data.(string)+ctx.Value(testContextKey).(string))

			return err
		}, "encoder", "encoder-alias")
	})

	t.Run("encoder-without-context", func(t *testing.T) {
		testRegisterEncoder(t, func(writer io.Writer, data any) error {
			_, err := io.WriteString(writer, data.(string)+testContextValue)

			return err
		}, "encoder-without-context", "encoder-without-context-alias")
	})
}

func testRegisterEncoder[E coder.EncoderConstraint](t *testing.T, encoder E, name, alias string) {
	t.Helper()

	coder.RegisterEncoder(encoder, name, alias)

	for _, key := range []string{name, alias} {
		enc := coder.GetEncoder(key)
		if enc == nil {
			t.Error("encoder not found")
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
