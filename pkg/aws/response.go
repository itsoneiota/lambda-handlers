package aws

import (
	"fmt"
	"log/slog"
	"net/http"
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
		slog.Debug(fmt.Sprintf("%v", r))
		res, err := h(Context{r.RequestContext}, NewAWSRequest(r))

		return NewEvent(res), err
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
