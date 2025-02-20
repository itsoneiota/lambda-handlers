package aws

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
)

type LambdaCallback = func(request *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

func Start(
	hf *handler.Handler,
	opts ...Setter,
) {
	h := &Handler{
		handler: hf,
		Opt:     &Opt{},
	}

	for _, o := range opts {
		o(h.Opt)
	}

	lambda.Start(handle(h))
}

func handle(h *Handler) LambdaCallback {
	return func(r *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		ctx := NewAWSContext(r.RequestContext)
		req := NewAWSRequest(r)

		var result *handler.Response

		func() {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("Recovered from panic in handler",
						slog.Any("error", rec),
						slog.String("stacktrace", string(debug.Stack())),
					)

					result = &handler.Response{
						StatusCode: http.StatusInternalServerError,
						Headers:    h.handler.Headers,
						Body:       `{"error": "Internal Server Error"}`,
					}
				}
			}()

			if result == nil {
				result = h.handler.Function(ctx, req)
			}
		}()

		if is2XXRange(result.StatusCode) {
			for _, i := range h.interceptors() {
				result = i(ctx, req, result)
			}
		}

		return NewEvent(result), nil
	}
}

func NewEvent(r *handler.Response) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: r.StatusCode,
		Headers:    encodeHeaders(r.Headers),
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

func is2XXRange(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}
