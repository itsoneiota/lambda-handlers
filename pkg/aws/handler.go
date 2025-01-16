package aws

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
)

type LambdaCallback = func(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

type Middleware func(*http.Request) error
type Interceptor func(*ResponseWriter) error

type Handler struct {
	function http.HandlerFunc
	*Opt
}

func Start(
	hf http.HandlerFunc,
	opts ...Setter,
) {
	h := &Handler{
		function: hf,
		Opt:      &Opt{},
	}

	for _, o := range opts {
		o(h.Opt)
	}

	lambda.Start(
		handle(h),
	)
}

func handle(h *Handler) LambdaCallback {
	return func(r *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		resp := NewResponseWriter(h.headers())
		req, err := NewHttpRequest(r)
		if err != nil {
			return nil, err
		}

		vars := map[string]string{}
		for key, value := range r.PathParameters {
			vars[key] = value
		}
		req = mux.SetURLVars(req, vars)

		for _, middleware := range h.middlewares() {
			if err := middleware(req); err != nil {
				fmt.Println(err)
			}
		}

		h.function(resp, req)

		if !isOkRange(resp.StatusCode) {
			return NewEvent(resp), nil
		}

		for _, interceptor := range h.interceptors() {
			if err := interceptor(resp); err != nil {
				fmt.Println(err)
			}
		}

		return NewEvent(resp), nil
	}
}

func NewEvent(r *ResponseWriter) *events.APIGatewayProxyResponse {
	headers := map[string]string{}
	for k, v := range r.Header() {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: r.StatusCode,
		Headers:    headers,
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
