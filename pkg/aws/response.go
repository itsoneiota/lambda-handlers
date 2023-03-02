package aws

import (
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
		res, err := h(NewAWSRequest(r))

		return NewEvent(res), err
	}
}

func NewEvent(r *handler.Response) *events.APIGatewayProxyResponse {
	headers := map[string]string{}
	for k, v := range r.Headers {
		headers[k] = v[0]
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: r.StatusCode,
		Headers:    headers,
		Body:       r.Body,
	}
}
