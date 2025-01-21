package mux

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/stretchr/testify/assert"
)

type metasyntactic struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
	Baz string `json:"baz"`
}

func TestRun(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		assert.NoError(t, err)

		return &handler.Response{StatusCode: http.StatusOK, Body: string(b), Headers: http.Header{
			"foo": {
				"bar",
			},
		}}
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{}
	New(testHandler, WithHeaders(http.Header{
		"default": {"header"},
	})).Run()(resp, req)

	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"foo":"foo","bar":"","baz":""}`, string(resp.body))
	assert.NotEmpty(t, resp.headers)
	assert.Equal(t, "header", resp.headers.Get("default"))
	assert.Equal(t, "bar", resp.headers.Get("foo"))
}

func TestMiddlewares(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		m := &metasyntactic{
			Foo: fmt.Sprintf("%d", ctx.Value("foo").(int)),
			Bar: req.QueryParams().Get("bar"),
		}

		b, err := json.Marshal(m)
		assert.NoError(t, err)

		return &handler.Response{StatusCode: http.StatusOK, Body: string(b)}
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, WithMiddlewares(
		func(next handler.HandlerFunc) handler.HandlerFunc {
			return func(ctx handler.Contexter, req handler.Requester) *handler.Response {
				ctx.SetValue("foo", 1)
				return next(ctx, req)
			}
		},
		func(next handler.HandlerFunc) handler.HandlerFunc {
			return func(ctx handler.Contexter, req handler.Requester) *handler.Response {
				req.SetQueryByName("bar", "2")
				return next(ctx, req)
			}
		},
	)).Run()(resp, req)

	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"foo":"1","bar":"2","baz":""}`, string(resp.body))
}

func TestInterceptors(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		assert.NoError(t, err)

		return &handler.Response{StatusCode: http.StatusOK, Body: string(b)}
	}

	resp := &fakeResponseWriter{
		headers: http.Header{},
	}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, WithInterceptors(
		func(
			_ handler.Contexter,
			_ handler.Requester,
			resp *handler.Response,
		) *handler.Response {
			m := &metasyntactic{}
			err := json.Unmarshal([]byte(resp.Body), m)
			assert.NoError(t, err)

			m.Bar = "4"

			b, err := json.Marshal(m)
			assert.NoError(t, err)

			resp.Body = string(b)

			return resp
		},
	)).Run()(resp, req)

	assert.Equal(t, http.StatusOK, resp.statusCode)
	assert.Equal(t, `{"foo":"foo","bar":"4","baz":""}`, string(resp.body))
}
