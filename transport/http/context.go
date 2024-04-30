package http

import (
	"context"
	"net/http"
)

type contextKey int

const (
	contextRequestKey contextKey = iota
)

func WithContextRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, contextRequestKey, r)
}

func RequestFromContext(ctx context.Context) *http.Request {
	return ctx.Value(contextRequestKey).(*http.Request)
}
