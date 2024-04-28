package http

import (
	"context"
	"net/http"

	"github.com/kamilov/go-kit/endpoint"
)

type (
	Middleware func(http.Handler) http.Handler
)

func handler[Input, Output any](server *Server, endpoint endpoint.Endpoint[Input, Output]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			input  Input
			output Output
			err    error
		)
	}
}

func withHttpMiddlewares(controller http.Handler, middlewares ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range middlewares {
			controller = m(controller)
		}

		controller.ServeHTTP(w, r)
	}
}

func withEndpointMiddlewares[Input, Output any](
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) endpoint.Endpoint[Input, Output] {
	return func(ctx context.Context, input Input) (Output, error) {
		for _, middleware := range middlewares {
			controller = middleware(controller)
		}

		return controller(ctx, input)
	}
}
