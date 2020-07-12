package goutils

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetRequestVarString func
func GetRequestVarString(r *http.Request, key string) string {
	vars := mux.Vars(r)
	return vars[key]
}

// GetRequestVarInt func
func GetRequestVarInt(r *http.Request, key string) (int, error) {
	return ParseInt(GetRequestVarString(r, key))
}

// GetRequestFormValueString func
func GetRequestFormValueString(r *http.Request, key string) string {
	return r.FormValue(key)
}

// GetRequestFormValueInt func
func GetRequestFormValueInt(r *http.Request, key string) (int, error) {
	return ParseInt(GetRequestFormValueString(r, key))
}

// ParseInt func
func ParseInt(value string) (int, error) {
	return strconv.Atoi(value)
}
