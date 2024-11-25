package mux

import (
	"net/http"
	"time"
)

// Context is the aws request context.
type Context struct {
	*http.Request
	values map[string]any
}

// SourceIP returns the source ip that has made the request.
func (c Context) SourceIP() string {
	return c.RemoteAddr
}

// UnixNow returns the Epoch-formatted request time, in milliseconds.
func (c Context) UnixNow() int64 {
	return time.Now().UnixMilli()
}

// HttpMethod returns the http method of a request
func (c Context) HttpMethod() string {
	return c.Request.Method
}

// Stage returns the stage of the environment.
func (c Context) Stage() string {
	return "dev" // We will default this to "dev" for the mux server as it is for local development.
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
