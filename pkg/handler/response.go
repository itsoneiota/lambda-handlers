package handler

import (
	"encoding/json"
	"net/http"
)

// Genertic Handler object which is the reciever in every handler method
type ResponseHandler struct {
	DefaultHeaders map[string]string
	responder      Responder
	logger         Logger
}

func NewResponseHandler(logger Logger, responder Responder) *ResponseHandler {
	return &ResponseHandler{
		logger:    logger,
		responder: responder,
	}
}

// Body gets request payload
func (r *ResponseHandler) BuildResponse(code int, model interface{}) (Responder, error) {
	body := ""
	if model != nil {
		bodyBytes, err := json.Marshal(model)
		if err != nil {
			r.responder.SetStatusCode(http.StatusInternalServerError)

			return r.responder, err
		}

		body = string(bodyBytes)
	}

	return r.BuildResponder(code, body)
}

// BuildRawJSONResponse builds an Response with the given status code & response body
// The Response will contain the raw response body and appropriate JSON header
func (r *ResponseHandler) BuildResponder(code int, body string) (Responder, error) {
	r.responder.SetStatusCode(code)
	r.responder.SetHeaders(r.DefaultHeaders)
	r.responder.SetBody(string(body))

	return r.responder, nil
}

func (r *ResponseHandler) BuildErrorResponse(err error) (Responder, error) {
	statusCode := http.StatusInternalServerError

	type statusCoder interface {
		StatusCode() int
	}

	servErr, hasStatusCode := err.(statusCoder)
	if hasStatusCode {
		statusCode = servErr.StatusCode()
	}

	if statusCode == http.StatusInternalServerError {
		r.logger.Error(err)
	}

	return r.BuildResponse(statusCode, err)
}