package log

import (
	"context"
	"time"

	"github.com/kamilov/go-kit/log"
	"github.com/loghole/dbhook"
)

const durationCtxKey = "duration"

type Hook struct {
	logger log.Logger
}

func New(logger log.Logger) *Hook {
	return &Hook{logger}
}

func (h *Hook) Before(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	return context.WithValue(ctx, durationCtxKey, time.Now()), nil
}

func (h *Hook) After(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.logger.
		With("Caller", input.Caller).
		Debugf(
			"Query: `%s`, Args: `%q`. duration: %s",
			input.Query,
			input.Args,
			time.Since(ctx.Value(durationCtxKey).(time.Time)),
		)
	return ctx, nil
}

func (h *Hook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.logger.
		With("Caller", input.Caller, "Error", input.Error).
		Errorf(
			"Query: `%s`, Args: `%q`. duration: %s",
			input.Query,
			input.Args,
			time.Since(ctx.Value(durationCtxKey).(time.Time)),
		)
	return ctx, nil
}
