package example

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsoneiota/lambda-handlers/internal/mocks"
	"github.com/itsoneiota/lambda-handlers/pkg/aws"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestFind_AWS(t *testing.T) {
	expectToken := "authToken"
	model := ExampleModel{}
	expectQuery := "POSTCODE"
	awsReq := &events.APIGatewayProxyRequest{
		Path: "test/123",
		QueryStringParameters: map[string]string{
			"postcode": expectQuery,
		},
		Headers: map[string]string{
			"Accept":        "application/json",
			"Authorization": expectToken,
		},
	}
	req := aws.NewAWSRequest(awsReq)

	// Mocks
	c := new(mocks.Connector)
	c.On("Authorize",
		expectToken,
	).Return(
		nil,
	).Times(1)

	c.On("Find",
		expectQuery,
	).Return(
		model,
		nil,
	).Times(1)

	l := test.NewNullLogger()
	slog.SetDefault(l)

	headers := http.Header{
		"Content-Type": []string{"application/json"},
	}

	resHander := handler.New(
		headers,
	)

	// Asserts
	resp := resHander.Run(FindHandler(resHander, c, nil, nil))(aws.Context{}, req)

	awsRes := aws.NewEvent(resp)
	expectAwsRes := &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		Body:              "{\"success\":false}",
		IsBase64Encoded:   false,
	}

	assert.IsType(t, &events.APIGatewayProxyResponse{}, awsRes)
	assert.Equal(t, expectAwsRes, awsRes)

}
