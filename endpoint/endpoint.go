package endpoint

import "context"

type Endpoint[Input, Output any] func(context.Context, Input) (Output, error)
type Middleware[Input, Output any] func(Endpoint[Input, Output]) Endpoint[Input, Output]

func Chain[Input, Output any](middlewares ...Middleware[Input, Output]) Middleware[Input, Output] {
	return func(next Endpoint[Input, Output]) Endpoint[Input, Output] {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}
