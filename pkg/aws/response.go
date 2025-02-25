package aws

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/pkg/serviceerror"
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
					slog.Error("Handler error",
						slog.Any("error", rec),
						slog.Group("stacktrace", "lines", stackTrace(debug.Stack())),
					)

					e := serviceerror.InternalServerError("Internal Server Error")
					b, _ := json.Marshal(e)

					result = &handler.Response{
						StatusCode: http.StatusInternalServerError,
						Headers:    h.handler.Headers,
						Body:       string(b),
					}
				}
			}()

			if result == nil {
				result = h.handler.Function(ctx, req)
				headers := http.Header{}
				if result.Headers != nil {
					headers = result.Headers
				}

				for k, v := range h.handler.Headers {
					for _, val := range v {
						headers.Add(k, val)
					}
				}

				result.Headers = headers
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

func stackTrace(stack []byte) []string {
	lines := strings.Split(string(stack), "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
