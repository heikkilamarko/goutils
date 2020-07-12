package goutils

// DataResponse struct
type DataResponse struct {
	Meta interface{} `json:"meta,omitempty"`
	Data interface{} `json:"data"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Error *ErrorResponseError `json:"error"`
}

// ErrorResponseError struct
type ErrorResponseError struct {
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// NewDataResponse func
func NewDataResponse(data, meta interface{}) *DataResponse {
	return &DataResponse{meta, data}
}

// NewErrorResponse func
func NewErrorResponse(code string, details map[string]string) *ErrorResponse {
	return &ErrorResponse{
		Error: &ErrorResponseError{code, details},
	}
}

// NewBadRequestResponse func
func NewBadRequestResponse(details map[string]string) *ErrorResponse {
	return NewErrorResponse(ErrCodeBadRequest, details)
}

// NewUnauthorizedResponse func
func NewUnauthorizedResponse(details map[string]string) *ErrorResponse {
	return NewErrorResponse(ErrCodeUnauthorized, details)
}

// NewNotFoundResponse func
func NewNotFoundResponse(details map[string]string) *ErrorResponse {
	return NewErrorResponse(ErrCodeNotFound, details)
}

// NewInternalErrorResponse func
func NewInternalErrorResponse(details map[string]string) *ErrorResponse {
	return NewErrorResponse(ErrCodeInternalError, details)
}
