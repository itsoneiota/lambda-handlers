package aws

import (
	"encoding/json"
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
	bodyStr := string(body)
	if !isOkRange(w.StatusCode) && !isValidJSON(bodyStr) {
		e := NewServiceError(
			GetServiceErrorCode(w.StatusCode),
			GetServiceErrorCode(w.StatusCode),
			bodyStr,
		)

		b, err := json.Marshal(e)
		if err != nil {
			return 0, err
		}

		bodyStr = string(b)
	}

	w.Body = bodyStr

	return len(body), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func isValidJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}
