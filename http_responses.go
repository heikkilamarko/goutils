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
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Error ErrorResponseError `json:"error"`
}

type ErrorResponseError struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

func NewDataResponse(data, meta interface{}) *DataResponse {
	return &DataResponse{data, meta}
}

func NewErrorResponse(code string, details interface{}) *ErrorResponse {
	return &ErrorResponse{Error: ErrorResponseError{code, details}}
}

func WriteOK(w http.ResponseWriter, data, meta interface{}) {
	WriteResponse(w, http.StatusOK, NewDataResponse(data, meta))
}

func WriteCreated(w http.ResponseWriter, data, meta interface{}) {
	WriteResponse(w, http.StatusCreated, NewDataResponse(data, meta))
}

func WriteNoContent(w http.ResponseWriter) {
	WriteResponse(w, http.StatusNoContent, nil)
}

func WriteBadRequest(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusBadRequest, NewErrorResponse(errCodeBadRequest, details))
}

func WriteUnauthorized(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusUnauthorized, NewErrorResponse(errCodeUnauthorized, details))
}

func WriteNotFound(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusNotFound, NewErrorResponse(errCodeNotFound, details))
}

func WriteInternalError(w http.ResponseWriter, details interface{}) {
	WriteResponse(w, http.StatusInternalServerError, NewErrorResponse(errCodeInternalError, details))
}

func WriteValidationError(w http.ResponseWriter, err error) {
	var verr *ValidationError
	if errors.As(err, &verr) {
		WriteBadRequest(w, verr.ErrorMap)
	} else {
		WriteInternalError(w, nil)
	}
}

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
