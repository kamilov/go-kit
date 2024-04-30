package http

import (
	"errors"
	"net/http"
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
