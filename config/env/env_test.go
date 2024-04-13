//nolint:cyclop // normal cyclomatic
package env_test

import (
	"context"
	"os"
	"reflect"
	"slices"
	"testing"

	"github.com/kamilov/go-kit/config"
	"github.com/kamilov/go-kit/config/env"
)

type Slices struct {
	Chars []byte `env:"TEST_ENV_SLICE_CHARS"`
	Int   []int  `env:"TEST_ENV_SLICE_INT"`
}

func TestParser(t *testing.T) {
	vars := map[string]string{
		"TEST_ENV_STRING":      "test",
		"TEST_ENV_BOOL":        "TRUE",
		"TEST_ENV_INT":         "1000",
		"TEST_ENV_UINT":        "1000",
		"TEST_ENV_FLOAT":       "10.01",
		"TEST_ENV_SLICE_CHARS": "abc",
		"TEST_ENV_SLICE_INT":   "[1, 2, 3]",
		"TEST_ENV_MAP":         `{"key":"value"}`,
	}

	for key, value := range vars {
		_ = os.Setenv(key, value)
	}

	compareValues := func(
		String string,
		Bool bool,
		Int int64,
		Uint uint64,
		Float float64,
		Chars []byte,
		Ints []int,
		Map map[string]string,
	) bool {
		return String == "test" &&
			Bool &&
			Int == int64(1000) &&
			Uint == uint64(1000) &&
			Float == 10.01 &&
			slices.Equal(Chars, []byte{'a', 'b', 'c'}) &&
			slices.Equal(Ints, []int{1, 2, 3}) &&
			reflect.DeepEqual(Map, map[string]string{"key": "value"})
	}

	t.Run("test default env", func(t *testing.T) {
		t.Helper()

		var test = struct {
			String string  `env:"TEST_ENV_STRING"`
			Bool   bool    `env:"TEST_ENV_BOOL"`
			Int    int64   `env:"TEST_ENV_INT"`
			Uint   uint64  `env:"TEST_ENV_UINT"`
			Float  float64 `env:"TEST_ENV_FLOAT"`
			Slices
			Map map[string]string `env:"TEST_ENV_MAP"`
		}{}

		if err := config.Read(context.Background(), &test); err != nil {
			t.Error(err)
		} else if !compareValues(
			test.String,
			test.Bool,
			test.Int,
			test.Uint,
			test.Float,
			test.Chars,
			test.Slices.Int,
			test.Map,
		) {
			t.Error("failed to parse env")
		}
	})

	t.Run("test default env", func(t *testing.T) {
		t.Helper()

		var test = struct {
			String string            `e_n_v:"STRING"`
			Bool   bool              `e_n_v:"BOOL"`
			Int    int64             `e_n_v:"INT"`
			Uint   uint64            `e_n_v:"UINT"`
			Float  float64           `e_n_v:"FLOAT"`
			Chars  []byte            `e_n_v:"SLICE_CHARS"`
			Ints   []int             `e_n_v:"SLICE_INT"`
			Map    map[string]string `e_n_v:"MAP"`
		}{}
		ctx := context.Background()
		ctx = context.WithValue(ctx, env.TagName, "e_n_v")
		ctx = context.WithValue(ctx, env.Prefix, "TEST_ENV_")

		if err := config.Read(ctx, &test); err != nil {
			t.Error(err)
		} else if !compareValues(test.String, test.Bool, test.Int, test.Uint, test.Float, test.Chars, test.Ints, test.Map) {
			t.Error("failed to parse env")
		}
	})

	for key := range vars {
		_ = os.Unsetenv(key)
	}
}
