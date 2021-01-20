package goutils

import (
	"testing"
)

func TestNewValidationError(t *testing.T) {
	var _ error = NewValidationError(map[string]string{})
}

func TestNewValidationErrorEmpty(t *testing.T) {
	verr := NewValidationError(map[string]string{})

	l := len(verr.ValidationErrors)

	if l != 0 {
		t.Errorf("len(ValidationErrors) = %d; want 0", l)
	}
}

func TestNewValidationErrorNonEmpty(t *testing.T) {
	verr := NewValidationError(map[string]string{
		"id": "invalid id",
	})

	l := len(verr.ValidationErrors)

	if l != 1 {
		t.Errorf("len(ValidationErrors) = %d; want 1", l)
	}
}
