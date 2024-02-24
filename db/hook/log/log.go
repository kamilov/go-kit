// Package log query log hook
package log

import (
	"context"
	"time"

	"github.com/kamilov/go-kit/log"
	"github.com/loghole/dbhook"
)

type txKey int

const durationCtxKey txKey = iota

// Hook log struct of a hook
type Hook struct {
	logger log.Logger
}

// New logger hook constructor
func New(logger log.Logger) *Hook {
	return &Hook{logger}
}

// Before callback
func (h *Hook) Before(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	return context.WithValue(ctx, durationCtxKey, time.Now()), nil
}

// After callback
func (h *Hook) After(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.logger.
		With("Caller", input.Caller).
		Debugf(
			"Query: `%s`, Args: `%q`. duration: %s",
			input.Query,
			input.Args,
			time.Since(ctx.Value(durationCtxKey).(time.Time)), //nolint:forcetypeassert // force
		)

	return ctx, nil
}

// Error callback
func (h *Hook) Error(ctx context.Context, input *dbhook.HookInput) (context.Context, error) {
	h.logger.
		With("Caller", input.Caller, "Error", input.Error).
		Errorf(
			"Query: `%s`, Args: `%q`. duration: %s",
			input.Query,
			input.Args,
			time.Since(ctx.Value(durationCtxKey).(time.Time)), //nolint:forcetypeassert // force
		)

	return ctx, nil
}
