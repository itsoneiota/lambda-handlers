package handler

import (
	"mime/multipart"
	"net/http"
	"net/url"
)

// Generic Request object which is used in every handler
type Requester interface {
	Body() string
	GetAuthToken() string
	Headers() http.Header
	MultipartReader() (*multipart.Reader, error)
	PathByName(name string) string
	QueryByName(name string) string
	QueryParams() url.Values
	SetQueryByName(name, set string)
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
	HttpMethod() string
	SetValue(string, any)
	Stage() string
	Value(string) any
}

// BeforeHandlerHook is a callback function called before a handler functions main logic is ran.
// A Callback function can be passed in when building a handler and is passed the raw API Gateway Request struct
type BeforeHandlerHook func(Requester) error

type HandlerFunc = func(c Contexter, request Requester) (*Response, error)

func WithValue(ctx Contexter, key string, value any) Contexter {
	ctx.SetValue(key, value)

	return ctx
}
