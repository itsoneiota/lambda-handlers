package handler

import (
	"mime/multipart"
	"net/http"
	"net/url"
)

type fakeRequest struct{}

func (c fakeRequest) Body() string {
	return ""
}

func (c fakeRequest) GetAuthToken() string {
	return ""
}

func (c fakeRequest) Headers() http.Header {
	return nil
}

func (c fakeRequest) MultipartReader() (*multipart.Reader, error) {
	return nil, nil
}

func (c fakeRequest) PathByName(string) string {
	return ""
}

func (c fakeRequest) QueryByName(string) string {
	return ""
}

func (c fakeRequest) QueryParams() url.Values {
	return nil
}

func (c fakeRequest) SetQueryByName(name, set string) {}

type fakeContext struct {
	values map[string]any
}

func (c fakeContext) SourceIP() string {
	return ""
}

func (c fakeContext) UnixNow() int64 {
	return 0
}

func (c fakeContext) UserAgent() string {
	return ""
}

func (c fakeContext) HttpMethod() string {
	return ""
}

func (c fakeContext) SetValue(key string, value any) {
	c.values[key] = value
}

func (c fakeContext) Stage() string {
	return ""
}

func (c fakeContext) Value(key string) any {
	return c.values[key]
}
