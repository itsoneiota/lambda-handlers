package aws

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ResponseWriter struct {
	*events.APIGatewayProxyResponse
}

func NewResponseWriter(headers http.Header) *ResponseWriter {
	h := map[string]string{}
	for k, v := range headers {
		if len(v) > 0 {
			h[k] = v[0]
		}
	}

	return &ResponseWriter{
		APIGatewayProxyResponse: &events.APIGatewayProxyResponse{
			Headers: h,
		},
	}
}

func (w *ResponseWriter) Header() http.Header {
	result := http.Header{}
	for k, v := range w.Headers {
		result.Add(k, v)
	}

	return result
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	w.Body = string(body)

	return len(body), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}
