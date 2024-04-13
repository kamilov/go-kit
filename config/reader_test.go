package config_test

import (
	"context"
	"errors"
	"testing"

	_ "github.com/kamilov/go-kit/coder/json"
	"github.com/kamilov/go-kit/config"
)

type contextKey string

const (
	testContextKey   contextKey = "testContextKey"
	testContextValue string     = "testContextValue"
)

type (
	testReader struct {
		Key   contextKey
		Value string
	}

	testReadFile struct {
		A string `json:"a" yaml:"a"`
		B int    `json:"b" yaml:"b"`
		C struct {
			D bool `json:"d" yaml:"d"`
		} `json:"c" yaml:"c"`
	}
)

func testRegisterReader[T config.ReaderConstraint](t *testing.T, readers ...T) {
	t.Helper()

	for _, reader := range readers {
		config.RegisterReader(reader)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, testContextKey, testContextValue)
	cfg := testReader{}

	if err := config.Read(ctx, &cfg); err != nil {
		t.Error(err)
	} else if cfg.Key != testContextKey || cfg.Value != testContextValue {
		t.Error("should be read data")
	}
}

func TestRegisterReader(t *testing.T) {
	t.Run("read config with context", func(t *testing.T) {
		testRegisterReader(t, func(_ context.Context, data any) error {
			impl, ok := data.(*testReader)
			if !ok {
				return nil
			}

			impl.Key = testContextKey

			return nil
		}, func(ctx context.Context, data any) error {
			impl, ok := data.(*testReader)
			if !ok {
				return nil
			}

			if value, exist := ctx.Value(testContextKey).(string); exist {
				impl.Value = value
			}

			return nil
		})
	})

	t.Run("read config without context", func(t *testing.T) {
		testRegisterReader(t, func(data any) error {
			impl, ok := data.(*testReader)
			if !ok {
				return nil
			}

			impl.Key = testContextKey

			return nil
		}, func(data any) error {
			impl, ok := data.(*testReader)
			if !ok {
				return nil
			}

			impl.Value = testContextValue

			return nil
		})
	})

	t.Run("read config error", func(t *testing.T) {
		var expectError = errors.New("test error")

		config.RegisterReader(func(data any) error {
			switch d := data.(type) {
			case error:
				return d
			default:
				return nil
			}
		})

		err := config.Read(context.Background(), expectError)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestReadFile(t *testing.T) {
	t.Run("read json file", func(t *testing.T) {
		var test testReadFile

		if err := config.ReadFile(context.Background(), "./stub/config.json", &test); err != nil {
			t.Error(err)
		} else if !testReadFileCheckData(test) {
			t.Errorf("error reading data from file: %v", test)
		}
	})

	t.Run("read yaml file", func(t *testing.T) {
		var test testReadFile

		if err := config.ReadFile(context.Background(), "./stub/config.yaml", &test); err == nil {
			t.Error("expected error as file decoder not exists")
		}
	})

	t.Run("read not be pointer", func(t *testing.T) {
		var test testReadFile

		if err := config.ReadFile(context.Background(), "./stub/bad.json", test); err == nil {
			t.Error("expected error as decode file")
		}
	})

	t.Run("read bad json file", func(t *testing.T) {
		var test testReadFile

		if err := config.ReadFile(context.Background(), "./stub/bad.json", &test); err == nil {
			t.Error("expected error as decode file")
		}
	})

	t.Run("read not found file", func(t *testing.T) {
		var test testReadFile

		if err := config.ReadFile(context.Background(), "./stub/config.xxx", &test); err == nil {
			t.Error("expected error")
		}
	})
}

func testReadFileCheckData(data testReadFile) bool {
	return data.A == "test" && data.B == 100 && data.C.D
}
