package aws

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaCallback = func(request *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse

func Start(
	h http.HandlerFunc,
	defaultHeaders http.Header,
) {
	lambda.Start(
		getHandler(h, defaultHeaders),
	)
}

func getHandler(
	h http.HandlerFunc,
	defaultHeaders http.Header,
) LambdaCallback {
	return func(r *events.APIGatewayProxyRequest) *events.APIGatewayProxyResponse {
		resp := NewResponseWriter(defaultHeaders)
		h(resp, NewHttpRequest())

		return NewEvent(resp)
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
