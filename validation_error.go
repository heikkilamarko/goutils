package goutils

import "encoding/json"

type ValidationError struct {
	ErrorMap map[string]string
}

func NewValidationError(errorMap map[string]string) *ValidationError {
	return &ValidationError{errorMap}
}

func (v *ValidationError) Error() string {
	message, err := json.Marshal(v.ErrorMap)
	if err != nil {
		return ""
	}
	return string(message)
}
