package yaml_test

import (
	"context"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

type test struct {
	A string `yaml:"a-key"`
	B int    `yaml:"b-key"`
}

//nolint:gochecknoglobals // a global description of the options for the test and the benchmark
var options = coder.TestOptions[test]{
	Name:    "yaml",
	Encoded: "a-key: test\nb-key: 100",
	Decoded: test{A: "test", B: 100},
}

func TestYaml(t *testing.T) {
	coder.Test(t, context.Background(), options)
}

func BenchmarkYaml(b *testing.B) {
	coder.Benchmark(b, context.Background(), options)
}
