package http

import (
	"bytes"
	"encoding"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
)

type (
	StatusError int

	Error struct {
		err  error
		code int
	}
)

const (
	ErrNotAcceptable        = StatusError(http.StatusNotAcceptable)
	ErrUnsupportedMediaType = StatusError(http.StatusUnsupportedMediaType)
)

func (e StatusError) Error() string {
	return http.StatusText(e.StatusCode())
}

func (e StatusError) StatusCode() int {
	return int(e)
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) StatusCode() int {
	if e.code != 0 {
		return e.code
	}

	var impl StatusCoder
	if errors.As(e.err, &impl) {
		return impl.StatusCode()
	}

	return http.StatusInternalServerError
}

func (e Error) Is(err error) bool {
	return errors.Is(e.err, err) || errors.Is(StatusError(e.code), err)
}

func (e Error) MarshalText() ([]byte, error) {
	var impl encoding.TextMarshaler
	if errors.As(e.err, &impl) {
		return impl.MarshalText()
	}

	return []byte(e.Error()), nil
}

func (e Error) MarshalJSON() ([]byte, error) {
	var impl json.Marshaler
	if errors.As(e.err, &impl) {
		return impl.MarshalJSON()
	}

	var buf bytes.Buffer

	buf.WriteString(`{"message": `)
	buf.WriteString(strconv.Quote(e.Error()))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e Error) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	var impl xml.Marshaler
	if errors.As(e.err, &impl) {
		return impl.MarshalXML(enc, start)
	}

	start = xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}
