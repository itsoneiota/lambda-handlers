package aws

import (
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
)

type AWSRequest struct {
	body        string
	pathParams  map[string]string
	queryParams url.Values
	headers     http.Header
}

func NewAWSRequest(r *events.APIGatewayProxyRequest) *AWSRequest {
	headers := http.Header{}
	for k, v := range r.Headers {
		headers.Set(k, v)
	}

	values := url.Values{}
	for k, v := range r.QueryStringParameters {
		values.Set(k, v)
	}

	return &AWSRequest{
		body:        r.Body,
		pathParams:  r.PathParameters,
		queryParams: values,
		headers:     headers,
	}
}

// Body gets request payload
func (r *AWSRequest) Body() string {
	return r.body
}

// Headers get the request headers
func (r *AWSRequest) Headers() http.Header {
	return r.headers
}

// PathByName gets a path parameter by its name eg. "productID"
func (r *AWSRequest) PathByName(name string) string {
	return r.pathParams[name]
}

// QueryByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) QueryByName(name string) string {
	return r.queryParams.Get(name)
}

// QueryByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) QueryParams() url.Values {
	return r.queryParams
}

// PathByName sets a query parameter by its name eg. "locale"
// This is used to alter requests in middleware functions.
func (r *AWSRequest) SetQueryByName(name, set string) {
	r.queryParams.Set(name, set)
}

// PathByName gets a query parameter by its name eg. "locale"
func (r *AWSRequest) GetAuthToken() string {
	if r.Headers().Get("Authorization") != "" {
		return r.Headers().Get("Authorization")
	} else {
		return r.Headers().Get("authorization")
	}
}
