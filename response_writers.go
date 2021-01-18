package goutils

import (
	"encoding/json"
	"errors"
	"net/http"
)

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
func WriteBadRequest(w http.ResponseWriter, details map[string]string) {
	WriteResponse(w, http.StatusBadRequest, NewBadRequestResponse(details))
}

// WriteUnauthorized writes 401 response
func WriteUnauthorized(w http.ResponseWriter, details map[string]string) {
	WriteResponse(w, http.StatusUnauthorized, NewUnauthorizedResponse(details))
}

// WriteNotFound writes 404 response
func WriteNotFound(w http.ResponseWriter, details map[string]string) {
	WriteResponse(w, http.StatusNotFound, NewNotFoundResponse(details))
}

// WriteInternalError writes 500 response
func WriteInternalError(w http.ResponseWriter, details map[string]string) {
	WriteResponse(w, http.StatusInternalServerError, NewInternalErrorResponse(details))
}

// WriteValidationError writes 400 or 500 response
func WriteValidationError(w http.ResponseWriter, err error) {
	var verr *ValidationError
	if errors.As(err, &verr) {
		WriteBadRequest(w, verr.ValidationErrors)
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
