package toml

import (
	"io"

	"github.com/kamilov/go-kit/coder"
	"github.com/pelletier/go-toml"
)

//nolint:gochecknoinits // used for automatic adding encode and decode functions
func init() {
	name := "toml"
	aliases := []string{"tml", "application/toml"}

	coder.RegisterDecoder(decoder, name, aliases...)
	coder.RegisterEncoder(encoder, name, aliases...)
}

func decoder(reader io.Reader, data any) error {
	return toml.NewDecoder(reader).Decode(data)
}

func encoder(writer io.Writer, data any) error {
	return toml.NewEncoder(writer).Encode(data)
}
