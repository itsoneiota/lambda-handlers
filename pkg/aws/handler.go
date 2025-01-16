package aws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
)

type LambdaCallback = func(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

type Middleware func(*http.Request) (*http.Request, error)
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
			var err error
			req, err = middleware(req)
			if err != nil {
				return errorResponse(resp, serviceerror.NewFromErr(err, ""))
			}
		}

		h.function(resp, req)

		if !isOkRange(resp.StatusCode) {
			return NewEvent(resp)
		}

		for _, interceptor := range h.interceptors() {
			if err := interceptor(resp); err != nil {
				return errorResponse(resp, serviceerror.NewFromErr(err, ""))
			}
		}

		return NewEvent(resp)
	}
}

func NewEvent(w *ResponseWriter) (*events.APIGatewayProxyResponse, error) {
	headers := map[string]string{}
	for k, v := range w.Header() {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: w.StatusCode,
		Headers:    headers,
		Body:       w.Body,
	}, nil
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

func errorResponse(w *ResponseWriter, srvErr *serviceerror.ServiceError) (*events.APIGatewayProxyResponse, error) {
	w.WriteHeader(srvErr.StatusCode())

	b, err := json.Marshal(srvErr)
	if err != nil {
		slog.Error(err.Error())
		b, _ = json.Marshal(serviceerror.InternalServerError(err.Error()))
	}

	w.Write(b)

	return NewEvent(w)
}
