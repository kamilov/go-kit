// Package xml provides encoding and decoding of XML data.
package xml

import (
	"encoding/xml"
	"io"

	"github.com/kamilov/go-kit/coder"
)

func init() {
	name := "xml"
	aliases := []string{"application/xml", "text/xml"}

	coder.RegisterDecoder(decode, name, aliases...)
	coder.RegisterEncoder(encode, name, aliases...)
}

func decode(reader io.Reader, data any) error {
	return xml.NewDecoder(reader).Decode(data)
}

func encode(writer io.Writer, data any) error {
	return xml.NewEncoder(writer).Encode(data)
}
