package goji

import (
	"net/http"
)

type dispatch struct{}

func (d dispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if h := HandlerFromContext(ctx); h != nil {
		h.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}
