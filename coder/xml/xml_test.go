package xml_test

import (
	"context"
	"testing"

	"github.com/kamilov/go-kit/coder"
)

type test struct {
	A string `xml:"a"`
	B int    `xml:"b"`
}

//nolint:gochecknoglobals // a global description of the options for the test and the benchmark
var options = coder.TestOptions[test]{
	Name:    "xml",
	Encoded: "<test><a>test</a><b>100</b></test>",
	Decoded: test{A: "test", B: 100},
}

func TestXML(t *testing.T) {
	coder.Test(t, context.Background(), options)
}

func BenchmarkXML(b *testing.B) {
	coder.Benchmark(b, context.Background(), options)
}
