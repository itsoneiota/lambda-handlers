package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Genertic Handler object which is the reciever in every handler method
type ResponseHandler struct {
	res http.ResponseWriter
}

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

// BuildResponseWithHeader creates an output Response with header
func (r *ResponseHandler) BuildResponseWithHeader(
	res http.ResponseWriter,
	code int,
	model interface{},
	headers http.Header,
) error {
	body := ""
	if model != nil {
		bodyBytes, err := json.Marshal(model)
		if err != nil {
			return err
		}

		body = string(bodyBytes)
	}

	return r.BuildResponderWithHeader(res, code, body, headers)
}

// BuildResponse creates an output Response
func (r *ResponseHandler) BuildResponse(
	res http.ResponseWriter,
	code int,
	model interface{},
) error {
	return r.BuildResponseWithHeader(res, code, model, http.Header{})
}

// BuildResponderWithHeader builds an Response with the given status code & response body
// The Response will contain the raw response body and appropriate JSON header
func (r *ResponseHandler) BuildResponderWithHeader(
	res http.ResponseWriter,
	code int,
	body string,
	inputHeaders http.Header,
) error {
	for k, vals := range inputHeaders {
		for _, v := range vals {
			res.Header().Add(k, v)
		}
	}

	res.Write([]byte(body))
	res.WriteHeader(code)

	return nil
}

// BuildResponder builds an Response with the given status code & response body
// The Response will contain the raw response body and appropriate JSON header
func (r *ResponseHandler) BuildResponder(
	res http.ResponseWriter,
	code int,
	body string,
) error {
	return r.BuildResponderWithHeader(res, code, body, http.Header{})
}

func (r *ResponseHandler) BuildErrorResponse(
	res http.ResponseWriter,
	err error,
) error {
	return r.BuildErrorResponseWithHeader(res, err, http.Header{})
}

func (r *ResponseHandler) BuildErrorResponseWithHeader(
	res http.ResponseWriter,
	err error,
	headers http.Header,
) error {
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

	if statusCode == http.StatusInternalServerError {
		slog.Error(err.Error())
	}

	return r.BuildResponseWithHeader(res, statusCode, serviceErr, headers)
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
