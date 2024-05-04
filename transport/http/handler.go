package http

import (
	"net/http"

	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/endpoint"
	"github.com/kamilov/go-kit/transport/http/content"
)

type (
	StatusCoder interface {
		StatusCode() int
	}

	Headerer interface {
		Header() http.Header
	}
)

func handler[Input, Output any](
	server *Server,
	controller endpoint.Endpoint[Input, Output],
	chain endpoint.Middleware[Input, Output],
) http.HandlerFunc {
	decode := newRequestDecoder[Input]()
	isNil := NewNilCheck(*new(Output))

	return func(w http.ResponseWriter, r *http.Request) {
		contentType := content.NegotiateContentType(r, server.negotiateTypes...)

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
			response, err = chain(controller)(ctx, input)
			if err != nil {
				response = Error{err, http.StatusInternalServerError}
			}
		}

		w.Header().Set("Content-Type", contentType+"; charset=utf-8")

		if impl, ok := response.(Headerer); ok {
			for key, values := range impl.Header() {
				for _, value := range values {
					w.Header().Add(key, value)
				}
			}
		}

		if impl, ok := response.(StatusCoder); ok {
			w.WriteHeader(impl.StatusCode())
		}

		if err == nil && isNil(response) {
			response = nil
			w.Header().Set("Content-Length", "0")
			w.WriteHeader(http.StatusNoContent)
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
