package goutils

import (
	"errors"
)

const (
	// ErrCodeBadRequest constant
	ErrCodeBadRequest = "bad_request"

	// ErrCodeNotFound constant
	ErrCodeNotFound = "not_found"

	// ErrCodeUnauthorized constant
	ErrCodeUnauthorized = "unauthorized"

	// ErrCodeInternalError constant
	ErrCodeInternalError = "internal_error"
)

var (
	// ErrNotFound error
	ErrNotFound = errors.New(ErrCodeNotFound)

	// ErrInternalError error
	ErrInternalError = errors.New(ErrCodeInternalError)
)
