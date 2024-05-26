package endpoint

import "context"

type Query[Input, Output any] interface {
	Handle(context.Context, Input) (Output, error)
}
