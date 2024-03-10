// Package yaml provides encoding and decoding of YAML data.
package yaml

import (
	"io"

	"github.com/kamilov/go-kit/coder"
	"gopkg.in/yaml.v3"
)

func init() {
	name := "yaml"
	aliases := []string{"yml", "text/yaml", "application/x-yaml", "text/x-yaml", "text/vnd.yaml"}

	coder.RegisterDecoder(decode, name, aliases...)
	coder.RegisterEncoder(encode, name, aliases...)
}

func decode(reader io.Reader, data any) error {
	return yaml.NewDecoder(reader).Decode(data)
}

func encode(writer io.Writer, data any) error {
	return yaml.NewEncoder(writer).Encode(data)
}
