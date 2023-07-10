package aws

import (
	"github.com/aws/aws-lambda-go/events"
)

// Context is the aws request context.
type Context struct {
	events.APIGatewayProxyRequestContext
}

// SourceIP returns the source ip that has made the request.
func (c Context) SourceIP() string {
	return c.Identity.SourceIP
}

// UnixNow returns the Epoch-formatted request time, in milliseconds.
func (c Context) UnixNow() int64 {
	return c.RequestTimeEpoch
}

// UserAgent returns the clients User-Agent request header value.
func (c Context) UserAgent() string {
	return c.Identity.UserAgent
}

// HttpMethod returns the http method that has been request.
func (c Context) HttpMethod() string {
	return c.HTTPMethod
}
