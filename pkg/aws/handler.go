package aws

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type LambdaCallback = func(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func Start(
	h http.HandlerFunc,
	beforeHook handler.BeforeHandlerHook,
	afterHook handler.AfterHandlerHook,
	defaultHeaders http.Header,
) {
	lambda.Start(
		getHandler(h, beforeHook, afterHook, defaultHeaders),
	)
}

func getHandler(
	h http.HandlerFunc,
	beforeHook handler.BeforeHandlerHook,
	afterHook handler.AfterHandlerHook,
	defaultHeaders http.Header,
) LambdaCallback {
	return func(r *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		resp := NewResponseWriter(defaultHeaders)
		req, err := NewHttpRequest(context.Background(), r)
		if err != nil {
			return nil, err
		}

		cnt := true
		if beforeHook != nil {
			cnt = beforeHook(resp, req)
		}

		if cnt {
			h(resp, req)

			if (resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices) && afterHook != nil {
				afterHook(resp)
			}
		}

		return NewEvent(resp), nil
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
