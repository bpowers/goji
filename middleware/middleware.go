/*
Package middleware contains utilities for Goji Middleware authors.

Unless you are writing middleware for your application, you should avoid
importing this package. Instead, use the abstractions provided by your
middleware package.
*/
package middleware

import (
	"context"
	"net/http"

	"goji.io"
	"goji.io/internal"
)

/*
Pattern returns the most recently matched Pattern, or nil if no pattern was
matched.
*/
var Pattern = goji.PatternFromContext

/*
SetPattern returns a new context in which the given Pattern is used as the most
recently matched pattern.
*/
func SetPattern(ctx context.Context, p goji.Pattern) context.Context {
	return context.WithValue(ctx, internal.Pattern, p)
}

/*
Handler returns the handler corresponding to the most recently matched Pattern,
or nil if no pattern was matched.

The handler returned by this function is the one that will be dispatched to at
the end of the middleware stack. If the returned Handler is nil, http.NotFound
will be used instead.
*/
var Handler = goji.HandlerFromContext

/*
SetHandler returns a new context in which the given Handler was most recently
matched and which consequently will be dispatched to.
*/
func SetHandler(ctx context.Context, h http.Handler) context.Context {
	return context.WithValue(ctx, internal.Handler, h)
}
