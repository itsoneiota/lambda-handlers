package aws

import (
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/handler"
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
		req, err := NewHttpRequest(r)
		if err != nil {
			return nil, err
		}

		vars := map[string]string{}
		for key, value := range r.PathParameters {
			vars[key] = value
		}
		req = mux.SetURLVars(req, vars)

		cnt := true
		if beforeHook != nil {
			cnt = beforeHook(resp, req)
		}

		if cnt {
			h(resp, req)

			if isOkRange(resp.StatusCode) && afterHook != nil {
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

func isOkRange(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}
