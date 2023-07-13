package handler

import (
	"net/http"
	"net/url"
)

// Generic Request object which is used in every handler
type Requester interface {
	Body() string
	Headers() http.Header
	PathByName(name string) string
	QueryByName(name string) string
	QueryParams() url.Values
	SetQueryByName(name, set string)
	GetAuthToken() string
}

// Generic Response object which is used in every handler
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       string
}

type Contexter interface {
	SourceIP() string
	UnixNow() int64
	UserAgent() string
}

type Logger interface {
	Error(args ...interface{})
}

// BeforeHandlerHook is a callback function called before a handler functions main logic is ran.
// A Callback function can be passed in when building a handler and is passed the raw API Gateway Request struct
type BeforeHandlerHook func(Requester) error

type HandlerFunc = func(c Contexter, request Requester) (*Response, error)
