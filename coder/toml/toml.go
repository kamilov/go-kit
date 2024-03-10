// Package toml provides encoding and decoding of TOML data.
package toml

import (
	"io"

	"github.com/kamilov/go-kit/coder"
	"github.com/pelletier/go-toml"
)

func init() {
	name := "toml"
	alias := "application/toml"

	coder.RegisterDecoder(decode, name, alias)
	coder.RegisterEncoder(encode, name, alias)
}

func decode(reader io.Reader, data any) error {
	return toml.NewDecoder(reader).Decode(data)
}

func encode(writer io.Writer, data any) error {
	return toml.NewEncoder(writer).Encode(data)
}
