package json_test

import (
	"context"
	"testing"

	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/coder/json"
)

type test struct {
	A string `json:"a"`
	B int    `json:"b"`
}

//nolint:gochecknoglobals // a global description of the options for the test and the benchmark
var options = coder.TestOptions[test]{
	Name:    "json",
	Encoded: `{"a":"test","b":100}`,
	Decoded: test{A: "test", B: 100},
}

func TestJSON(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, json.AllowUnknownFields, false)
	ctx = context.WithValue(ctx, json.UseNumber, true)
	ctx = context.WithValue(ctx, json.EscapeHTML, true)

	coder.Test(t, ctx, options)
	coder.Test(t, context.Background(), options)
}

func BenchmarkJSON(b *testing.B) {
	coder.Benchmark(b, context.Background(), options)
}
