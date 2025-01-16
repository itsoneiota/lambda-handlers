package example

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsoneiota/lambda-handlers/v2/internal/mocks"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/aws"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
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

func (s *FindHandlerSuite) Connector(err error) Connector {
	c := new(mocks.Connector)

	c.On("Authorize",
		s.token,
	).Return(
		err,
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
	res := aws.NewResponseWriter(
		http.Header{
			"Content-Type": {"application/json"},
		},
	)

	FindHandler(s.Connector(nil))(res, s.req)

	awsRes, err := aws.NewEvent(res)
	s.NoError(err)
	expectAwsRes := &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		Body:              `{"success":false}`,
		IsBase64Encoded:   false,
	}

	s.IsType(&events.APIGatewayProxyResponse{}, awsRes)
	s.Equal(expectAwsRes, awsRes)
}

func (s *FindHandlerSuite) TestHandlerError() {
	res := aws.NewResponseWriter(
		http.Header{
			"Content-Type": {"application/json"},
		},
	)

	FindHandler(s.Connector(serviceerror.BadRequest("something bad has happened")))(res, s.req)

	awsRes, err := aws.NewEvent(res)
	s.NoError(err)
	expectAwsRes := &events.APIGatewayProxyResponse{
		StatusCode:        http.StatusBadRequest,
		Headers:           map[string]string{"Content-Type": "application/json"},
		MultiValueHeaders: map[string][]string(nil),
		Body:              `{"error":{"id":"BAD_REQUEST","code":"BAD_REQUEST","message":"something bad has happened"}}`,
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
