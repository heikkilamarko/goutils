package goutils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errCodeBadRequest    = http.StatusText(http.StatusBadRequest)
	errCodeUnauthorized  = http.StatusText(http.StatusUnauthorized)
	errCodeNotFound      = http.StatusText(http.StatusNotFound)
	errCodeInternalError = http.StatusText(http.StatusInternalServerError)
)

// DataResponse struct
type DataResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}

// ErrorResponseError struct
type ErrorResponseError struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// NewDataResponse func
func NewDataResponse(data, meta interface{}) *DataResponse {
	return &DataResponse{data, meta}
}

// NewErrorResponse func
func NewErrorResponse(code string, details interface{}) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponseError{code, details}}
}

// WriteOK writes 200 response
func WriteOK(w http.ResponseWriter, data, meta interface{}) {
	WriteResponse(w, http.StatusOK, NewDataResponse(data, meta))
}

// WriteCreated writes 201 response
func WriteCreated(w http.ResponseWriter, data, meta interface{}) {
	WriteResponse(w, http.StatusCreated, NewDataResponse(data, meta))
}

// WriteNoContent writes 204 response
func WriteNoContent(w http.ResponseWriter) {
	WriteResponse(w, http.StatusNoContent, nil)
}

// WriteBadRequest writes 400 response
func WriteBadRequest(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusBadRequest, NewErrorResponse(errCodeBadRequest, details))
}

// WriteUnauthorized writes 401 response
func WriteUnauthorized(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusUnauthorized, NewErrorResponse(errCodeUnauthorized, details))
}

// WriteNotFound writes 404 response
func WriteNotFound(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusNotFound, NewErrorResponse(errCodeNotFound, details))
}

// WriteInternalError writes 500 response
func WriteInternalError(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusInternalServerError, NewErrorResponse(errCodeInternalError, details))
}

// WriteValidationError writes 400 or 500 response
func WriteValidationError(w http.ResponseWriter, err error) {
	var verr *ValidationError
	if errors.As(err, &verr) {
		WriteBadRequest(w, verr.ErrorMap)
	} else {
		WriteInternalError(w, nil)
	}
}

// WriteResponse func
func WriteResponse(w http.ResponseWriter, code int, body interface{}) {
	if body != nil {
		content, err := json.Marshal(body)

		if err != nil {
			WriteInternalError(w, nil)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(content)
	} else {
		w.WriteHeader(code)
	}
}
