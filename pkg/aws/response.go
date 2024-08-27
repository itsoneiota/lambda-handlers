package aws

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type LambdaCallback = func(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func Start(
	h handler.HandlerFunc,
	defaultHeaders http.Header,
) {
	lambda.Start(
		getHandler(h, defaultHeaders),
	)
}

func getHandler(
	h handler.HandlerFunc,
	defaultHeaders http.Header,
) LambdaCallback {
	return func(r *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		resp := NewResponseWriter(defaultHeaders)
		err := h(resp, NewAWSRequest(r))

		return NewEvent(resp), err
	}
}

func NewEvent(r *ResponseWriter) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: r.StatusCode,
		Headers:    r.Headers,
		Body:       r.Body,
	}
}

func encodeHeaders(h http.Header) map[string]string {
	result := map[string]string{}

	for hKey := range h {
		valsUnique := unique(h.Values(hKey))
		result[hKey] = strings.Join(valsUnique, "; ")
	}

	return result
}

func unique(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, value := range slice {
		if !encountered[value] {
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}
