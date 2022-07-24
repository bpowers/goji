package goji

import (
	"context"
	"net/http"

	"goji.io/internal"
)

// PatternFromContext returns the most recently matched Pattern, or nil if no pattern was
// matched.
func PatternFromContext(ctx context.Context) Pattern {
	if pi := ctx.Value(internal.Pattern); pi != nil {
		if p, ok := pi.(Pattern); ok {
			return p
		}
	}
	return nil
}

// HandlerFromContext returns the handler corresponding to the most recently matched Pattern,
// or nil if no pattern was matched.
//
// The handler returned by this function is the one that will be dispatched to at
// the end of the middleware stack. If the returned Handler is nil, http.NotFound
// will be used instead.
func HandlerFromContext(ctx context.Context) http.Handler {
	if hi := ctx.Value(internal.Handler); hi != nil {
		if h, ok := hi.(http.Handler); ok {
			return h
		}
	}
	return nil
}
