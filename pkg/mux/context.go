package mux

import (
	"net/http"
	"time"
)

// Context is the aws request context.
type Context struct {
	*http.Request
}

// SourceIP returns the source ip that has made the request.
func (c Context) SourceIP() string {
	// TODO: functionally test this
	return c.RemoteAddr
}

// UnixNow returns the Epoch-formatted request time, in milliseconds.
func (c Context) UnixNow() int64 {
	// TODO: functionally test this
	return time.Now().UnixMilli()
}

// UnixNow returns the Epoch-formatted request time, in milliseconds.
func (c Context) HTTPMethod() string {
	// TODO: functionally test this
	return c.Request.Method
}
