package handler

import (
	"mime/multipart"
	"net/http"
	"net/url"
)

// Generic Request object which is used in every handler
type Requester interface {
	AddCookie(c *http.Cookie)
	Body() string
	Context() Contexter
	Cookie(name string) (*http.Cookie, error)
	Cookies() []*http.Cookie
	GetAuthToken() string
	Headers() http.Header
	MultipartReader() (*multipart.Reader, error)
	PathByName(name string) string
	QueryByName(name string) string
	QueryParams() url.Values
	Referer() string
	SetQueryByName(name, set string)
	UserAgent() string
}

type Contexter interface {
	SourceIP() string
	UnixNow() int64
	UserAgent() string
	HttpMethod() string
	Stage() string
}

// BeforeHandlerHook is a callback function called before a handler functions main logic is ran.
// A Callback function can be passed in when building a handler and is passed the raw API Gateway Request struct
type BeforeHandlerHook func(Requester) error

type HandlerFunc = func(res http.ResponseWriter, req Requester) error
