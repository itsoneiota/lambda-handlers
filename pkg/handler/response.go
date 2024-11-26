package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Setter func(*opt)

type opt struct {
	headers http.Header
}

// Generic Response object which is used in every handler
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

func NewResponse(
	statusCode int,
	body any,
	opts ...Setter,
) *Response {
	opt := &opt{}
	for _, o := range opts {
		o(opt)
	}

	var b string
	if v, ok := body.(string); ok {
		b = v
	} else {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return &Response{StatusCode: http.StatusInternalServerError}
		}

		b = string(bodyBytes)
	}

	return &Response{
		StatusCode: statusCode,
		Body:       b,
		Headers:    opt.headers,
	}
}

func NewErrorResponse(
	err error,
	opts ...Setter,
) *Response {
	statusCode := http.StatusInternalServerError
	var serviceErr error

	isServiceError, code := isServiceError(err)

	if isServiceError {
		statusCode = code
		serviceErr = err
	} else {
		// If its a general error - we don't want to return the message as its a code/integration issue.
		// We don't want those messages being shown to users.
		serviceErr = &ServiceError{
			Err: Error{
				ID:      "UNKNOWN_ERROR",
				Code:    "UNKNOWN_ERROR",
				Message: "An unknown error occurred",
			},
		}
	}

	slog.Error(err.Error())

	return NewResponse(statusCode, serviceErr, opts...)
}

func isServiceError(err error) (bool, int) {
	var code int

	type serviceError interface {
		Code() string
		Error() string
		StatusCode() int
	}

	se, isSe := err.(serviceError)

	if isSe {
		code = se.StatusCode()
	}

	return isSe, code
}

func WithHeaders(headers http.Header) Setter {
	return func(o *opt) {
		o.headers = headers
	}
}
