package api

import (
	"context"
	"crypto/rand"
)

type requestIDKey struct{}

// WithRequestID returns a context with a new local log correlation ID.
// Callers that log before Do can attach the ID here so their logs and
// the client's retry/error logs refer to the same logical request.
// It is not sent to the server and is not a deduplication key.
func WithRequestID(ctx context.Context) context.Context {
	return context.WithValue(ctx, requestIDKey{}, rand.Text())
}

// RequestID returns the local log correlation ID, or "" if none is attached.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}
