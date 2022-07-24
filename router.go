package goji

import (
	"context"
	"net/http"

	"goji.io/internal"
)

type match struct {
	context.Context
	p Pattern
	h http.Handler
}

func (m match) Value(key interface{}) interface{} {
	switch key.(type) {
	case internal.PatternContextKey:
		return m.p
	case internal.HandlerContextKey:
		return m.h
	default:
		return m.Context.Value(key)
	}
}

var _ context.Context = match{}
