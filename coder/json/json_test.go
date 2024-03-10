package json_test

import (
	"context"
	"testing"

	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/coder/json"
)

type test struct {
	A string `json:"a-key"`
	B int    `json:"b-key"`
}

var options = coder.Options[test]{
	Name:    "json",
	Encoded: `{"a-key":"test","b-key":100}`,
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
