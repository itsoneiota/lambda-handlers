package example

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsoneiota/lambda-handlers/internal/mocks"
	"github.com/itsoneiota/lambda-handlers/pkg/aws"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/stretchr/testify/suite"
)

type FindHandlerSuite struct {
	suite.Suite
	token string
	resp  ExampleModel
	query string
	req   *http.Request
}

func (s *FindHandlerSuite) SetupTest() {
	s.token = "authToken"
	s.resp = ExampleModel{}
	s.query = "POSTCODE"
	s.req = &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path:     "find",
			RawQuery: fmt.Sprintf("postcode=%s", s.query),
		},
		Header: http.Header{
			"Accept":        {"application/json"},
			"Authorization": {s.token},
		},
	}

}

func (s *FindHandlerSuite) Connector() Connector {
	c := new(mocks.Connector)

	c.On("Authorize",
		s.token,
	).Return(
		nil,
	).Times(1)

	c.On("Find",
		s.query,
	).Return(
		s.resp,
		nil,
	).Times(1)

	return c
}

func (s *FindHandlerSuite) TestHandler() {
	resHander := handler.NewResponseHandler()
	res := aws.NewResponseWriter(
		http.Header{
			"Content-Type": {"application/json"},
		},
	)

	// Asserts
	FindHandler(resHander, s.Connector(), nil, nil)(res, s.req)

	awsRes := aws.NewEvent(res)
	expectAwsRes := &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		Body:              "{\"success\":false}",
		IsBase64Encoded:   false,
	}

	s.IsType(&events.APIGatewayProxyResponse{}, awsRes)
	s.Equal(expectAwsRes, awsRes)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFindHandlerSuite(t *testing.T) {
	suite.Run(t, new(FindHandlerSuite))
}
