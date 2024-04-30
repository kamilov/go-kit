package text

import (
	"errors"
	"strconv"

	"github.com/kamilov/go-kit/coder"
)

const (
	numDefaultSize = strconv.IntSize
	num8BitSize    = 8
	num16BitSize   = 16
	num32BitSize   = 32
	num64BitSize   = 64
	numBase        = 10
)

var (
	ErrUnsupportedType = errors.New("unsupported content type")
)

//nolint:gochecknoinits // use init to register encoders and decoders func
func init() {
	name := "text"
	aliases := []string{"text/plain", "text/html"}

	coder.RegisterDecoder(decoder, name, aliases...)
	coder.RegisterEncoder(encoder, name, aliases...)
}
