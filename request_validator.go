package goutils

import "net/http"

// RequestValidator struct
type RequestValidator struct {
	Request          *http.Request
	ValidationErrors map[string]string
}

// IsValid method
func (v *RequestValidator) IsValid() bool {
	return len(v.ValidationErrors) == 0
}
