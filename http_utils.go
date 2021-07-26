package goutils

import (
	"net/http"
	"strings"
)

// TokenFromHeader func
func TokenFromHeader(r *http.Request) string {
	a := r.Header.Get("Authorization")
	if 7 < len(a) && strings.ToUpper(a[0:6]) == "BEARER" {
		return a[7:]
	}
	return ""
}
