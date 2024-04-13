package yaml

import (
	"io"

	"github.com/kamilov/go-kit/coder"
	"gopkg.in/yaml.v3"
)

//nolint:gochecknoinits // used for automatic adding encode and decode functions
func init() {
	name := "yaml"
	aliases := []string{"yml", "text/yaml", "application/x-yaml", "text/x-yaml", "text/vnd.yaml"}

	coder.RegisterDecoder(decoder, name, aliases...)
	coder.RegisterEncoder(encoder, name, aliases...)
}

func decoder(reader io.Reader, data any) error {
	return yaml.NewDecoder(reader).Decode(data)
}

func encoder(writer io.Writer, data any) error {
	return yaml.NewEncoder(writer).Encode(data)
}
