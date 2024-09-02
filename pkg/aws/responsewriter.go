package aws

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsoneiota/lambda-handlers/pkg/serviceerror"
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
	if !isOkRange(w.StatusCode) && !isValidJSONObject(bodyStr) {
		var decodedString string
		if err := json.Unmarshal([]byte(bodyStr), &decodedString); err == nil {
			bodyStr = decodedString
		}

		e := serviceerror.NewServiceError(
			serviceerror.GetServiceErrorCode(w.StatusCode),
			serviceerror.GetServiceErrorCode(w.StatusCode),
			bodyStr,
		)

		b, err := json.Marshal(e)
		if err != nil {
			slog.Error(err.Error())
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

func isValidJSONObject(s string) bool {
	var js interface{}
	err := json.Unmarshal([]byte(s), &js)
	if err != nil {
		return false
	}

	switch js.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}
