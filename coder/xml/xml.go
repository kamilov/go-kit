package xml

import (
	"encoding/xml"
	"io"

	"github.com/kamilov/go-kit/coder"
)

//nolint:gochecknoinits // used for automatic adding encode and decode functions
func init() {
	name := "xml"
	aliases := []string{"application/xml", "text/xml"}

	coder.RegisterDecoder(decoder, name, aliases...)
	coder.RegisterEncoder(encoder, name, aliases...)
}

func decoder(reader io.Reader, data any) error {
	return xml.NewDecoder(reader).Decode(data)
}

func encoder(writer io.Writer, data any) error {
	return xml.NewEncoder(writer).Encode(data)
}
