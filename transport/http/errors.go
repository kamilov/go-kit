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

	HTTPError struct {
		err  error
		code int
	}
)

const (
	ErrNotAcceptable        = StatusError(http.StatusNotAcceptable)
	ErrUnsupportedMediaType = StatusError(http.StatusUnsupportedMediaType)
)

func Error(err error, code int) HTTPError {
	return HTTPError{err, code}
}

func (e StatusError) Error() string {
	return http.StatusText(e.StatusCode())
}

func (e StatusError) StatusCode() int {
	return int(e)
}

func (e HTTPError) Error() string {
	return e.err.Error()
}

func (e HTTPError) Unwrap() error {
	return e.err
}

func (e HTTPError) StatusCode() int {
	if e.code != 0 {
		return e.code
	}

	var impl StatusCoder
	if errors.As(e.err, &impl) {
		return impl.StatusCode()
	}

	return http.StatusInternalServerError
}

func (e HTTPError) Is(err error) bool {
	return errors.Is(e.err, err) || errors.Is(StatusError(e.code), err)
}

func (e HTTPError) MarshalText() ([]byte, error) {
	var impl encoding.TextMarshaler
	if errors.As(e.err, &impl) {
		return impl.MarshalText()
	}

	return []byte(e.Error()), nil
}

func (e HTTPError) MarshalJSON() ([]byte, error) {
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

func (e HTTPError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	var impl xml.Marshaler
	if errors.As(e.err, &impl) {
		return impl.MarshalXML(enc, start)
	}

	start = xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}
