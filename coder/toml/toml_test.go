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

var options = coder.Options[test]{
	Name:    "toml",
	Encoded: "a = \"test\"\nb = 100",
	Decoded: test{A: "test", B: 100},
}

func TestYaml(t *testing.T) {
	coder.Test(t, context.Background(), options)
}

func BenchmarkYaml(b *testing.B) {
	coder.Benchmark(b, context.Background(), options)
}
