package aws

import (
	"github.com/aws/aws-lambda-go/events"
)

// Context is the aws request context.
type Context struct {
	events.APIGatewayProxyRequestContext
	values map[string]any
}

func NewAWSContext(ctx events.APIGatewayProxyRequestContext) *Context {
	return &Context{
		APIGatewayProxyRequestContext: ctx,
		values:                        map[string]any{},
	}
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

func (c Context) Stage() string {
	return c.APIGatewayProxyRequestContext.Stage
}

func (c Context) Value(key string) any {
	if val, ok := c.values[key]; ok {
		return val
	}

	return nil
}

func (c Context) SetValue(key string, value any) {
	c.values[key] = value
}
