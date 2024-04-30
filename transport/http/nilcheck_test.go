package http_test

import (
	"testing"

	"github.com/kamilov/go-kit/transport/http"
)

func TestNilCheck(t *testing.T) {
	tests := []struct {
		message  string
		expected bool
		zero     any
		data     any
	}{
		{
			"nil interface",
			true,
			nil,
			nil,
		},
		{
			"empty string",
			true,
			"",
			"",
		},
		{
			"nil struct",
			true,
			(*struct{})(nil),
			(*struct{})(nil),
		},
		{
			"zero struct",
			false,
			struct{}{},
			struct{}{},
		},
		{
			"nil map",
			true,
			(map[string]string)(nil),
			(map[string]string)(nil),
		},
		{
			"zero map",
			false,
			(map[string]string)(nil),
			map[string]string{},
		},
		{
			"non-zero map",
			false,
			(map[string]string)(nil),
			map[string]string{"foo": "bar"},
		},
		{
			"nil slice",
			true,
			([]string)(nil),
			([]string)(nil),
		},
		{
			"zero slice",
			false,
			([]string)(nil),
			[]string{},
		},
		{
			"non-zero slice",
			false,
			([]string)(nil),
			[]string{"a", "b", "c"},
		},
		{
			"non-zero slice pointer",
			false,
			([]string)(nil),
			func() any {
				m := []string{"a", "b", "c"}
				return &m
			}(),
		},
		{
			"boolean",
			false,
			false,
			false,
		},
		{
			"integer",
			false,
			0,
			0,
		},
	}

	for _, test := range tests {
		isNil := http.NewNilCheck(test.zero)

		if isNil(test.data) != test.expected {
			t.Errorf("%s should be %t", test.message, test.expected)
		}
	}
}
