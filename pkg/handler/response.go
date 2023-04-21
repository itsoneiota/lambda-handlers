package handler

import (
	"encoding/json"
	"net/http"
)

// Genertic Handler object which is the reciever in every handler method
type ResponseHandler struct {
	defaultHeaders http.Header
	logger         Logger
}

func NewResponseHandler(
	logger Logger,
	defaultHeads http.Header,
) *ResponseHandler {
	return &ResponseHandler{
		logger:         logger,
		defaultHeaders: defaultHeads,
	}
}

// Body gets request payload
func (r *ResponseHandler) BuildResponse(code int, model interface{}, headers http.Header) (*Response, error) {
	body := ""
	if model != nil {
		bodyBytes, err := json.Marshal(model)
		if err != nil {
			return &Response{StatusCode: http.StatusInternalServerError}, err
		}

		body = string(bodyBytes)
	}

	return r.BuildResponder(code, body, headers)
}

// BuildRawJSONResponse builds an Response with the given status code & response body
// The Response will contain the raw response body and appropriate JSON header
func (r *ResponseHandler) BuildResponder(code int, body string, inputHeaders http.Header) (*Response, error) {
	// add default and input headers
	h := r.defaultHeaders
	for inputKey, inputVals := range inputHeaders {
		for _, v := range inputVals {
			h.Add(inputKey, v)
		}
	}

	return &Response{
		StatusCode: code,
		Body:       body,
		Headers:    h,
	}, nil
}

func (r *ResponseHandler) BuildErrorResponse(err error, headers http.Header) (*Response, error) {
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
		r.logger.Error(err)
	}

	return r.BuildResponse(statusCode, serviceErr, headers)
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
