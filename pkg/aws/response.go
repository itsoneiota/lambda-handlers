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

func Start(h handler.HandlerFunc) {
	lambda.Start(
		getHandler(h),
	)
}

func getHandler(h handler.HandlerFunc) LambdaCallback {
	return func(r *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		var result *handler.Response
		var err error

		func() {
			defer func() {
				if rec := recover(); rec != nil {
					slog.Error("Recovered from panic in handler",
						slog.Any("error", rec),
						slog.String("stacktrace", formatStackTrace(debug.Stack())),
					)

					result = &handler.Response{
						StatusCode: http.StatusInternalServerError,
						Headers:    http.Header{"Content-Type": []string{"application/json"}},
						Body:       `{"error": {"code":"INTERNAL_SERVER_ERROR", "id":"INTERNAL_SERVER_ERROR", "message":"Internal Server Error"}}`,
					}
				}
			}()

			if result == nil {
				result, err = h(Context{r.RequestContext}, NewAWSRequest(r))
				headers := http.Header{}
				if result.Headers != nil {
					headers = result.Headers
				}

				result.Headers = headers
			}
		}()

		return NewEvent(result), err
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

func formatStackTrace(stack []byte) string {
	lines := strings.Split(string(stack), "\n")
	for i, line := range lines {
		lines[i] = "    " + line // Indent for better readability
	}
	return "\n" + strings.Join(lines, "\n")
}
