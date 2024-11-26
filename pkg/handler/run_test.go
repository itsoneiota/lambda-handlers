package handler

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testHandler := func(ctx Contexter, req Requester) *Response {
		return &Response{StatusCode: http.StatusOK, Body: "foo"}
	}

	req := testRequest{}
	ctx := testContext{}
	handler := New(http.Header{})
	resp := handler.Run(testHandler)(ctx, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "foo", resp.Body)
}

func TestRunResponseError(t *testing.T) {
	testHandler := func(ctx Contexter, req Requester) *Response {
		return &Response{StatusCode: http.StatusNotFound, Body: "foo not found"}
	}

	req := testRequest{}
	ctx := testContext{}
	handler := New(http.Header{})
	resp := handler.Run(testHandler)(ctx, req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "foo not found", resp.Body)
}

type testRequest struct{}

func (c testRequest) Body() string {
	return ""
}

func (c testRequest) GetAuthToken() string {
	return ""
}

func (c testRequest) Headers() http.Header {
	return nil
}

func (c testRequest) MultipartReader() (*multipart.Reader, error) {
	return nil, nil
}

func (c testRequest) PathByName(string) string {
	return ""
}

func (c testRequest) QueryByName(string) string {
	return ""
}

func (c testRequest) QueryParams() url.Values {
	return nil
}

func (c testRequest) SetQueryByName(name, set string) {}

type testContext struct{}

func (c testContext) SourceIP() string {
	return ""
}

func (c testContext) UnixNow() int64 {
	return 0
}

func (c testContext) UserAgent() string {
	return ""
}

func (c testContext) HttpMethod() string {
	return ""
}

func (c testContext) SetValue(string, any) {}

func (c testContext) Stage() string {
	return ""
}

func (c testContext) Value(string) any {
	return nil
}
