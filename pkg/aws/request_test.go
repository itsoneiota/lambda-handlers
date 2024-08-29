package aws

import (
	"io"
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
			"Content-Type":  "application/json",
			"Authorization": "Bearer example-token",
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
	s.Equal(s.req.Path, req.URL.Path)
	s.Equal(s.req.Path, req.URL.Path)
	s.Equal("https", req.URL.Scheme)

	b, err := io.ReadAll(req.Body)
	s.NoError(err)
	s.Equal(s.req.Body, string(b))
}

// func TestRequestGetters(t *testing.T) {
// 	body := "abcdef"

// 	token := "abc.def.ghi"

// 	headName := "header"
// 	headVal := "h1"

// 	pathKey := "id"
// 	pathVal := "p1"

// 	path2Key := "subid"
// 	path2Val := "p2"

// 	queryKey := "q"
// 	queryVal := "football"

// 	query2Key := "t"
// 	query2Val := "red"

// 	headers := http.Header{}
// 	headers.Set(headName, headVal)
// 	headers.Set("Authorization", token)

// 	values := url.Values{}
// 	values.Set(queryKey, queryVal)
// 	values.Set(query2Key, query2Val)

// 	req := AWSRequest{
// 		body:    body,
// 		headers: headers,
// 		pathParams: map[string]string{
// 			pathKey:  pathVal,
// 			path2Key: path2Val,
// 		},
// 		queryParams: values,
// 	}

// 	assert.Equal(t, body, req.Body())

// 	assert.Equal(t, headVal, req.Headers().Get(headName))

// 	assert.Equal(t, token, req.Headers().Get("Authorization"))

// 	assert.Equal(t, pathVal, req.PathByName(pathKey))

// 	assert.Equal(t, queryVal, req.QueryByName(queryKey))
// }

// func TestBody_Empty(t *testing.T) {
// 	req := AWSRequest{}

// 	assert.Equal(t, "", req.Body())
// 	assert.Equal(t, "", req.Headers().Get("headName"))
// 	assert.Equal(t, "", req.Headers().Get("Authorization"))
// 	assert.Equal(t, "", req.PathByName("pathKey"))
// 	assert.Equal(t, "", req.QueryByName("queryKey"))
// }

// func TestSetQueryByName(t *testing.T) {
// 	queryKey := "q"
// 	queryVal := "football"

// 	query2Key := "t"
// 	query2Val := "red"

// 	newQueryVal := "soccer"

// 	values := url.Values{}
// 	values.Set(queryKey, queryVal)
// 	values.Set(query2Key, query2Val)

// 	req := AWSRequest{
// 		queryParams: values,
// 	}

// 	// BEFORE
// 	assert.Equal(t, queryVal, req.QueryByName(queryKey))
// 	assert.Equal(t, query2Val, req.QueryByName(query2Key))

// 	req.SetQueryByName(queryKey, newQueryVal)
// 	// AFTER
// 	assert.Equal(t, newQueryVal, req.QueryByName(queryKey))
// 	assert.Equal(t, query2Val, req.QueryByName(query2Key))
// }

// func TestNewAWSRequest(t *testing.T) {
// 	pathKey := "pathKey"
// 	pathVal := "path"
// 	queryKey := "q"
// 	queryVal := "football"
// 	body := "body here"
// 	headKey := "headerKey"
// 	headVal := "headerValue"

// 	req := &events.APIGatewayProxyRequest{
// 		QueryStringParameters: map[string]string{
// 			queryKey: queryVal,
// 		},
// 		PathParameters: map[string]string{
// 			pathKey: pathVal,
// 		},
// 		Body: body,
// 		Headers: map[string]string{
// 			headKey: headVal,
// 		},
// 	}

// 	actual := NewAWSRequest(req)

// 	// BEFORE
// 	assert.Equal(t, body, actual.Body())
// 	assert.Equal(t, pathVal, actual.PathByName(pathKey))
// 	assert.Equal(t, headVal, actual.Headers().Get(headKey))
// 	assert.Equal(t, queryVal, actual.QueryByName(queryKey))
// }

// func TestMultipartReader(t *testing.T) {
// 	headKey := "Content-Type"
// 	headVal := "multipart/form-data; boundary=BOUNDARY"
// 	content := "This is content"
// 	body := "--BOUNDARY\r\n" +
// 		"Content-Disposition: form-data; name=\"value\"\r\n" +
// 		"\r\n" +
// 		content +
// 		"\r\n--BOUNDARY--\r\n"

// 	req := &events.APIGatewayProxyRequest{
// 		Body: body,
// 		Headers: map[string]string{
// 			headKey: headVal,
// 		},
// 	}

// 	actual := NewAWSRequest(req)

// 	reader, err := actual.MultipartReader()
// 	assert.NoError(t, err)
// 	assert.IsType(t, &multipart.Reader{}, reader)

// 	part, err := reader.NextPart()
// 	assert.NoError(t, err)

// 	cnt, err := io.ReadAll(part)
// 	assert.NoError(t, err)
// 	assert.Equal(t, content, string(cnt))
// }

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}
