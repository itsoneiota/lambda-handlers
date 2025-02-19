package fakers

import (
	"mime/multipart"
	"net/http"
	"net/url"
)

type Request struct{}

func (c Request) Body() string {
	return ""
}

func (c Request) GetAuthToken() string {
	return ""
}

func (c Request) Headers() http.Header {
	return nil
}

func (c Request) MultipartReader() (*multipart.Reader, error) {
	return nil, nil
}

func (c Request) PathByName(string) string {
	return ""
}

func (c Request) QueryByName(...string) string {
	return ""
}

func (c Request) QueryParams() url.Values {
	return nil
}

func (c Request) SetQueryByName(name, set string) {}

type Context struct {
	Values map[string]any
}

func (c Context) SourceIP() string {
	return ""
}

func (c Context) UnixNow() int64 {
	return 0
}

func (c Context) UserAgent() string {
	return ""
}

func (c Context) HttpMethod() string {
	return ""
}

func (c Context) SetValue(key string, value any) {
	c.Values[key] = value
}

func (c Context) Stage() string {
	return ""
}

func (c Context) Value(key string) any {
	return c.Values[key]
}
