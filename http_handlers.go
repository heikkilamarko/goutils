package goutils

import "net/http"

// NotFoundHandler func
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteNotFound(w, nil)
	})
}
