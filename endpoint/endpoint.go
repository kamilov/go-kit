package endpoint

import "context"

type Endpoint[Input, Output any] func(context.Context, Input) (Output, error)
type Middleware[Input, Output any] func(Endpoint[Input, Output]) Endpoint[Input, Output]
