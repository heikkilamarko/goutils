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

type DataResponse struct {
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}

type ErrorResponseError struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

func NewDataResponse(data, meta any) *DataResponse {
	return &DataResponse{data, meta}
}

func NewErrorResponse(code string, details any) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponseError{code, details}}
}

func WriteOK(w http.ResponseWriter, data, meta any) {
	WriteResponse(w, http.StatusOK, NewDataResponse(data, meta))
}

func WriteCreated(w http.ResponseWriter, data, meta any) {
	WriteResponse(w, http.StatusCreated, NewDataResponse(data, meta))
}

func WriteNoContent(w http.ResponseWriter) {
	WriteResponse(w, http.StatusNoContent, nil)
}

func WriteBadRequest(w http.ResponseWriter, details any) {
	WriteResponse(w, http.StatusBadRequest, NewErrorResponse(errCodeBadRequest, details))
}

func WriteUnauthorized(w http.ResponseWriter, details any) {
	WriteResponse(w, http.StatusUnauthorized, NewErrorResponse(errCodeUnauthorized, details))
}

func WriteNotFound(w http.ResponseWriter, details any) {
	WriteResponse(w, http.StatusNotFound, NewErrorResponse(errCodeNotFound, details))
}

func WriteInternalError(w http.ResponseWriter, details any) {
	WriteResponse(w, http.StatusInternalServerError, NewErrorResponse(errCodeInternalError, details))
}

func WriteValidationError(w http.ResponseWriter, err error) {
	var verr *ValidationError
	if errors.As(err, &verr) {
		WriteBadRequest(w, verr.Errors)
	} else {
		WriteInternalError(w, nil)
	}
}

func WriteResponse(w http.ResponseWriter, code int, body any) {
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
