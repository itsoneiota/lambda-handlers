package serviceerror

import (
	"fmt"
	"net/http"
)

// Error codes
const (
	CodeUnknown             = "UNKNOWN_ERROR"
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeNotImplemented      = "NOT_IMPLEMENTED"
	CodeUnprocessableEntity = "UNPROCESSABLE_ENTITY"
	CodeConflict            = "CONFLICT"
	CodeRequestTimeout      = "REQUEST_TIMEOUT"
	CodeNotFound            = "NOT_FOUND"
	CodeForbidden           = "FORBIDDEN"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeBadRequest          = "BAD_REQUEST"
	CodeFound               = "FOUND"
	CodeMovedPermanently    = "MOVED_PERMANENTLY"
)

// StatusCodes mapped to the error codes
var StatusCodes = map[string]int{
	CodeInternalServerError: http.StatusInternalServerError,
	CodeNotImplemented:      http.StatusNotImplemented,
	CodeUnprocessableEntity: http.StatusUnprocessableEntity,
	CodeConflict:            http.StatusConflict,
	CodeRequestTimeout:      http.StatusRequestTimeout,
	CodeNotFound:            http.StatusNotFound,
	CodeForbidden:           http.StatusForbidden,
	CodeUnauthorized:        http.StatusUnauthorized,
	CodeBadRequest:          http.StatusBadRequest,
	CodeFound:               http.StatusFound,
	CodeMovedPermanently:    http.StatusMovedPermanently,
}

// defaultErrorMessages are default error messages if we are unable to get a message from the client error.
var defaultErrorMessages = map[string]string{
	CodeInternalServerError: "Internal Service Error",
	CodeNotImplemented:      "Not Implemented",
	CodeUnprocessableEntity: "Unprocessable Entity",
	CodeConflict:            "Conflict",
	CodeRequestTimeout:      "Request Timeout",
	CodeNotFound:            "Not Found",
	CodeForbidden:           "Forbidden",
	CodeUnauthorized:        "Unauthorized",
	CodeBadRequest:          "Bad Request",
	CodeFound:               "Found",
	CodeMovedPermanently:    "Moved Permanently",
}

// ServiceError - represents the service error
type ServiceError struct {
	Err Error `json:"error"`
}

// Error holds the error contents of the service error
type Error struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error returns the error message
func (se *ServiceError) Error() string {
	return se.Err.Message
}

// Code returns the error code
func (se *ServiceError) Code() string {
	return se.Err.Code
}

// StatusCode returns the errors StatusCode
func (se *ServiceError) StatusCode() int {
	respCode := http.StatusInternalServerError
	if val, ok := StatusCodes[se.Err.Code]; ok {
		respCode = val
	}
	return respCode
}

// New creates a new service error with the given code and message
func NewServiceError(id, code, message string) *ServiceError {
	if id == "" {
		id = code
	}

	return &ServiceError{
		Error{
			ID:      id,
			Code:    code,
			Message: message,
		},
	}
}

// NewFromErr returns a new service error built from an existing error
func NewFromErr(err error, message string) *ServiceError {
	code := CodeInternalServerError
	if e, ok := err.(*ServiceError); ok {
		code = e.Code()
	}
	return &ServiceError{
		Error{
			Code:    code,
			Message: fmt.Sprintf("%s: %s", message, err.Error()),
		},
	}
}

// InternalServerError is a helper method for creating a service error with an 'InternalServerError' code
func InternalServerError(message string) *ServiceError {
	return NewServiceError(CodeInternalServerError, CodeInternalServerError, message)
}

// NotImplemented is a helper method for creating a service error with an 'NotImplemented' code
func NotImplemented(message string) *ServiceError {
	return NewServiceError(CodeNotImplemented, CodeNotImplemented, message)
}

// UnprocessableEntity is a helper method for creating a service error with an 'UnprocessableEntity' code
func UnprocessableEntity(message string) *ServiceError {
	return NewServiceError(CodeUnprocessableEntity, CodeUnprocessableEntity, message)
}

// Conflict is a helper method for creating a service error with an 'Conflict' code
func Conflict(message string) *ServiceError {
	return NewServiceError(CodeConflict, CodeConflict, message)
}

// RequestTimeout is a helper method for creating a service error with an 'RequestTimeout' code
func RequestTimeout(message string) *ServiceError {
	return NewServiceError(CodeRequestTimeout, CodeRequestTimeout, message)
}

// NotFound is a helper method for creating a service error with an 'NotFound' code
func NotFound(message string) *ServiceError {
	return NewServiceError(CodeNotFound, CodeNotFound, message)
}

// Forbidden is a helper method for creating a service error with an 'Forbidden' code
func Forbidden(message string) *ServiceError {
	return NewServiceError(CodeForbidden, CodeForbidden, message)
}

// Unauthorized is a helper method for creating a service error with an 'Unauthorized' code
func Unauthorized(message string) *ServiceError {
	return NewServiceError(CodeUnauthorized, CodeUnauthorized, message)
}

// BadRequest is a helper method for creating a service error with an 'BadRequest' code
func BadRequest(message string) *ServiceError {
	return NewServiceError(CodeBadRequest, CodeBadRequest, message)
}

// Found is a helper method for creating a service error with an 'Found' code
func Found(message string) *ServiceError {
	return NewServiceError(CodeFound, CodeFound, message)
}

// MovedPermanently is a helper method for creating a service error with an 'MovedPermanently' code
func MovedPermanently(message string) *ServiceError {
	return NewServiceError(CodeMovedPermanently, CodeMovedPermanently, message)
}

// GetDefaultErrorMessage get the default error message
func GetDefaultErrorMessage(statusCode int) (string, int) {
	var message string
	var status int
	for k, v := range StatusCodes {
		if v == statusCode {
			message = defaultErrorMessages[k]
			status = StatusCodes[k]
		}
	}

	if message == "" {
		message = defaultErrorMessages[CodeInternalServerError]
		status = StatusCodes[CodeInternalServerError]
	}

	return message, status
}

// GetServiceErrorCode gets the service error code based on a given http status.
func GetServiceErrorCode(statusCode int) string {
	result := CodeInternalServerError
	for k, v := range StatusCodes {
		if v == statusCode {
			result = k
			break
		}
	}

	return result
}
