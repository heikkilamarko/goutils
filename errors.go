package goutils

import (
	"errors"
	"fmt"
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

// ValidationError error
type ValidationError struct {
	ValidationErrors map[string]string
}

// Error method
func (v *ValidationError) Error() string {
	return fmt.Sprintf("%v", v.ValidationErrors)
}

// NewValidationError func
func NewValidationError(validationErrors map[string]string) *ValidationError {
	return &ValidationError{validationErrors}
}
