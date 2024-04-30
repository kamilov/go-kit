package http

import (
	"context"
	"net/http"

	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/endpoint"
)

type (
	StatusCoder interface {
		StatusCode() int
	}

	Headerer interface {
		Header() http.Header
	}

	emptyBody struct{}
)

func EmptyBodyAdapter[Output any](fn func(ctx context.Context) (Output, error)) endpoint.Endpoint[emptyBody, Output] {
	return func(ctx context.Context, _ emptyBody) (Output, error) {
		return fn(ctx)
	}
}

func handler[Input, Output any](_ *Server, controller endpoint.Endpoint[Input, Output]) http.HandlerFunc {
	decode := newRequestDecoder[Input]()
	isNil := NewNilCheck(*new(Output))

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := NegotiateContentType(r)

		enc := coder.GetEncoder(contentType)
		if enc == nil {
			handleError(w, ErrNotAcceptable)
			return
		}

		var (
			input    Input
			response any
			err      error
		)

		ctx := WithContextRequest(r.Context(), r)
		err = decode(ctx, &input, r)
		if err != nil {
			response = Error{err, getStatusCode(err, http.StatusBadRequest)}
		} else {
			response, err = controller(ctx, input)
			if err != nil {
				response = Error{err, 0}
			}
		}

		w.Header().Set("Content-Type", contentType+"; charset=utf-8")

		if impl, ok := response.(Headerer); ok {
			for key, values := range impl.Header() {
				w.Header().Set(key, values[0])
			}
		}

		if impl, ok := response.(StatusCoder); ok {
			w.WriteHeader(impl.StatusCode())
		}

		if err == nil && isNil(response) {
			response = nil
			w.Header().Set("Content-Length", "0")
		}

		if err = enc(ctx, w, response); err != nil {
			handleError(w, err)
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	code := getStatusCode(err, http.StatusInternalServerError)
	if code < http.StatusInternalServerError {
		http.Error(w, err.Error(), code)
	} else {
		http.Error(w, http.StatusText(code), code)
	}
}

func getStatusCode(err error, fallback int) int {
	if impl, ok := err.(StatusCoder); ok {
		return impl.StatusCode()
	}
	return fallback
}
