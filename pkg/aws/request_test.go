package aws

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/suite"
)

type RequestSuite struct {
	suite.Suite
	req *events.APIGatewayProxyRequest
}

func (s *RequestSuite) SetupTest() {
	s.req = &events.APIGatewayProxyRequest{
		Resource:   "/products/{id}",
		Path:       "/products/ABC123",
		HTTPMethod: http.MethodPut,
		Headers: map[string]string{
			"Host":              "example.com",
			"X-Forwarded-Proto": "https",
			"Content-Type":      "application/json",
			"Authorization":     "Bearer example-token",
		},
		MultiValueHeaders: map[string][]string{
			"X-Custom-Header": {"value1", "value2"},
		},
		QueryStringParameters: map[string]string{
			"locale":   "en-GB",
			"currency": "GBP",
		},
		MultiValueQueryStringParameters: map[string][]string{
			"extend": {"attributes", "tabs"},
		},
		PathParameters: map[string]string{
			"id": "ABC123",
		},
		StageVariables: map[string]string{
			"env": "production",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			AccountID:  "1234567890",
			ResourceID: "resource-id",
			Stage:      "prod",
			RequestID:  "00000000-0000-0000-0000-000000000000",
			Identity: events.APIGatewayRequestIdentity{
				SourceIP:  "127.0.0.1",
				UserAgent: "Mozilla/5.0 (compatible; Example/0.1; +http://example.com)",
			},
			ResourcePath: "/products/{id}",
			HTTPMethod:   http.MethodPut,
			APIID:        "api-id",
		},
		Body:            "{\"name\": \"Example Product\"}",
		IsBase64Encoded: false,
	}
}

func (s *RequestSuite) TestNewHttpRequest() {
	req, err := NewHttpRequest(s.req)
	s.NoError(err)

	s.IsType(&http.Request{}, req)

	s.Equal(s.req.HTTPMethod, req.Method)

	s.IsType(&url.URL{}, req.URL)
	s.Equal("https", req.URL.Scheme)
	s.Equal("example.com", req.URL.Host)
	s.Equal("/products/ABC123", req.URL.Path)
	s.Equal("currency=GBP&extend=attributes&extend=tabs&locale=en-GB", req.URL.RawQuery)

	s.Equal(6, len(req.Header))

	b, err := io.ReadAll(req.Body)
	s.NoError(err)
	s.Equal(s.req.Body, string(b))

	s.Equal("example.com", req.Host)
	s.Equal("127.0.0.1", req.RemoteAddr)
}

func (s *RequestSuite) TestNewHttpRequestEncodedBody() {
	s.req.Body = "eyJuYW1lIjogIkV4YW1wbGUgUHJvZHVjdCJ9"
	s.req.IsBase64Encoded = true
	req, err := NewHttpRequest(s.req)
	s.NoError(err)

	s.IsType(&http.Request{}, req)

	b, err := io.ReadAll(req.Body)
	s.NoError(err)
	s.Equal("{\"name\": \"Example Product\"}", string(b))
}

func (s *RequestSuite) TestNewHttpRequestMultipartForm() {
	s.req.Headers = map[string]string{
		"Host":              "example.com",
		"X-Forwarded-Proto": "https",
		"Content-Type":      "multipart/form-data; boundary=BOUNDARY",
		"Authorization":     "Bearer example-token",
	}

	content := "This is content"
	s.req.Body = "--BOUNDARY\r\n" +
		"Content-Disposition: form-data; name=\"value\"\r\n" +
		"\r\n" +
		content +
		"\r\n--BOUNDARY--\r\n"

	req, err := NewHttpRequest(s.req)
	s.NoError(err)

	s.IsType(&http.Request{}, req)

	form := req.MultipartForm
	s.IsType(&multipart.Form{}, form)

	for _, vals := range form.Value {
		s.Equal(content, vals[0])
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}
