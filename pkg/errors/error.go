package errors

import (
	"fmt"
	"net/http"
)

const (
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
	CodeBadGateway          = "BAD_GATEWAY"
)

var statusCodes = map[string]int{
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
	CodeBadGateway:          http.StatusBadGateway,
}

var defaultErrorMessages = map[string]string{
	CodeInternalServerError: "Internal error",
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

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) StatusCode() int {
	if result, ok := statusCodes[e.Code]; ok {
		return result
	}

	return http.StatusInternalServerError
}

func (e *Error) Error() string {
	result := defaultErrorMessages[e.Code]
	if e.Message != "" {
		result = e.Message
	}

	return result
}

func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewFromErr(err error, message string) *Error {
	code := CodeInternalServerError
	if e, ok := err.(*Error); ok {
		code = e.Code
	}

	errMsg := err.Error()
	if message != "" {
		if err.Error() == "" {
			errMsg = fmt.Sprintf("%s", err.Error())
		} else {
			errMsg = fmt.Sprintf("%s: %s", message, err.Error())
		}
	}

	return New(code, errMsg)
}

// InternalServerError is a helper method for creating a error with an 'InternalServerError' code
func InternalServerError(message string) *Error {
	return New(CodeInternalServerError, message)
}

// NotImplemented is a helper method for creating a error with an 'NotImplemented' code
func NotImplemented(message string) *Error {
	return New(CodeNotImplemented, message)
}

// UnprocessableEntity is a helper method for creating a error with an 'UnprocessableEntity' code
func UnprocessableEntity(message string) *Error {
	return New(CodeUnprocessableEntity, message)
}

// Conflict is a helper method for creating a error with an 'Conflict' code
func Conflict(message string) *Error {
	return New(CodeConflict, message)
}

// RequestTimeout is a helper method for creating a error with an 'RequestTimeout' code
func RequestTimeout(message string) *Error {
	return New(CodeRequestTimeout, message)
}

// NotFound is a helper method for creating a error with an 'NotFound' code
func NotFound(message string) *Error {
	return New(CodeNotFound, message)
}

// Forbidden is a helper method for creating a error with an 'Forbidden' code
func Forbidden(message string) *Error {
	return New(CodeForbidden, message)
}

// Unauthorized is a helper method for creating a error with an 'Unauthorized' code
func Unauthorized(message string) *Error {
	return New(CodeUnauthorized, message)
}

// BadRequest is a helper method for creating a error with an 'BadRequest' code
func BadRequest(message string) *Error {
	return New(CodeBadRequest, message)
}

// Found is a helper method for creating a error with an 'Found' code
func Found(message string) *Error {
	return New(CodeFound, message)
}

// MovedPermanently is a helper method for creating a error with an 'MovedPermanently' code
func MovedPermanently(message string) *Error {
	return New(CodeMovedPermanently, message)
}

// BadGateway is a helper method for creating a error with an 'BadGateway' code
func BadGateway(message string) *Error {
	return New(CodeBadGateway, message)
}
