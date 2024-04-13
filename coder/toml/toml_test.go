package toml_test

import (
	"context"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

type test struct {
	A string `toml:"a"`
	B int    `toml:"b"`
}

//nolint:gochecknoglobals // a global description of the options for the test and the benchmark
var options = coder.TestOptions[test]{
	Name:    "toml",
	Encoded: "a = \"test\"\nb = 100",
	Decoded: test{A: "test", B: 100},
}

func TestTOML(t *testing.T) {
	coder.Test(t, context.Background(), options)
}

func BenchmarkTOML(b *testing.B) {
	coder.Benchmark(b, context.Background(), options)
}
