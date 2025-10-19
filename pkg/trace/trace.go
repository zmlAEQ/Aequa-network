package trace

import (
    "context"
)

type ctxKey struct{}

// WithTraceID returns a new context that carries the provided trace id.
func WithTraceID(ctx context.Context, id string) context.Context {
    if id == "" { return ctx }
    return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext extracts the trace id from context if present.
func FromContext(ctx context.Context) (string, bool) {
    v := ctx.Value(ctxKey{})
    if s, ok := v.(string); ok && s != "" { return s, true }
    return "", false
}