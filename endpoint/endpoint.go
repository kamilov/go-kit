package endpoint

import "context"

type (
	// Endpoint is the fundamental building block of servers and clients
	Endpoint[Input, Output any] func(context.Context, Input) (Output, error)
	// Middleware is a chainable behavior modifier for endpoints
	Middleware[Input, Output any] func(Endpoint[Input, Output]) Endpoint[Input, Output]
)

// Chain is a helper function for composing middlewares
func Chain[I, O any](middleware Middleware[I, O], middlewares ...Middleware[I, O]) Middleware[I, O] {
	return func(endpoint Endpoint[I, O]) Endpoint[I, O] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			endpoint = middlewares[i](endpoint)
		}
		return middleware(endpoint)
	}
}
