package endpoint

import "context"

type Command[Input any] interface {
	Handle(context.Context, Input) error
}
