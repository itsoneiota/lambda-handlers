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

type Contexter interface {
	SourceIP() string
	UnixNow() int64
	UserAgent() string
	HttpMethod() string
	SetValue(string, any)
	Stage() string
	Value(string) any
}

type HandlerFunc = func(c Contexter, request Requester) *Response

// Genertic Handler object which is the reciever in every handler method
type Handler struct {
	function HandlerFunc
	headers  http.Header
}

func New(
	function HandlerFunc,
	opts ...Setter,
) *Handler {
	opt := &opt{}
	for _, o := range opts {
		o(opt)
	}

	return &Handler{
		function: function,
		headers:  opt.headers,
	}
}

func WithValue(ctx Contexter, key string, value any) Contexter {
	ctx.SetValue(key, value)

	return ctx
}
