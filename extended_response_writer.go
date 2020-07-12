package goutils

import "net/http"

// ExtendedResponseWriter struct
type ExtendedResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader method
func (w *ExtendedResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// NewExtendedResponseWriter func
func NewExtendedResponseWriter(w http.ResponseWriter) *ExtendedResponseWriter {
	return &ExtendedResponseWriter{w, http.StatusOK}
}
