package aws

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/helpers"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
)

type ResponseWriter struct {
	*events.APIGatewayProxyResponse
	headers http.Header
}

func NewResponseWriter(headers http.Header) *ResponseWriter {
	if headers == nil {
		headers = http.Header{}
	}

	return &ResponseWriter{
		APIGatewayProxyResponse: &events.APIGatewayProxyResponse{},
		headers:                 headers,
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.headers
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	bodyStr := string(body)
	if !helpers.IsOkRange(w.StatusCode()) && !helpers.IsValidJSONObject(bodyStr) {
		var decodedString string
		if err := json.Unmarshal([]byte(bodyStr), &decodedString); err == nil {
			bodyStr = decodedString
		}

		e := serviceerror.NewServiceError(
			serviceerror.GetServiceErrorCode(w.StatusCode()),
			serviceerror.GetServiceErrorCode(w.StatusCode()),
			bodyStr,
		)

		b, err := json.Marshal(e)
		if err != nil {
			slog.Error(err.Error())
			return 0, err
		}

		bodyStr = string(b)
	} else if !helpers.IsValidJSONObject(bodyStr) {
		var decodedString string
		err := json.Unmarshal([]byte(bodyStr), &decodedString)
		if err == nil {
			bodyStr = decodedString
		}
	}

	w.APIGatewayProxyResponse.Body = bodyStr

	return len(body), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.APIGatewayProxyResponse.StatusCode = statusCode
}

func (w *ResponseWriter) Body() string {
	return w.APIGatewayProxyResponse.Body
}

func (w *ResponseWriter) StatusCode() int {
	return w.APIGatewayProxyResponse.StatusCode
}
