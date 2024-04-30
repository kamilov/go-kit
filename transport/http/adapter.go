package http

import (
	"context"
	"net/http"

	"github.com/kamilov/go-kit/endpoint"
)

type (
	empty    struct{}
	redirect struct {
		code int
		url  string
	}
)

func (r redirect) StatusCode() int {
	return r.code
}

func (r redirect) Header() http.Header {
	h := http.Header{}

	h.Set("Location", r.url)

	return h
}

func EmptyRequestAdapter[Output any](fn func(ctx context.Context) (Output, error)) endpoint.Endpoint[empty, Output] {
	return func(ctx context.Context, _ empty) (Output, error) {
		return fn(ctx)
	}
}

func EmptyResponseAdapter[Input any](fn func(ctx context.Context, input Input) error) endpoint.Endpoint[Input, *empty] {
	return func(ctx context.Context, input Input) (*empty, error) {
		err := fn(ctx, input)
		return nil, err
	}
}

func RedirectAdapter[Input any](url string, code int, fn func(ctx context.Context, input Input) error) endpoint.Endpoint[Input, *redirect] {
	return func(ctx context.Context, input Input) (*redirect, error) {
		err := fn(ctx, input)
		if err != nil {
			return nil, err
		}

		return &redirect{code, url}, err
	}
}
