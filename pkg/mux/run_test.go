package mux

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/itsoneiota/lambda-handlers/pkg/serviceerror"
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

	resp := &ResponseWriter{}
	req := &http.Request{}
	New(testHandler, WithHeaders(http.Header{
		"default": {"header"},
	})).Run()(resp, req)

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, `{"foo":"foo","bar":"","baz":""}`, string(resp.Body))
	assert.NotEmpty(t, resp.Headers)
	assert.Equal(t, "header", resp.Headers.Get("default"))
	assert.Equal(t, "bar", resp.Headers.Get("foo"))
}

func TestMiddlewares(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		m := &metasyntactic{
			Foo: "foo",
			Bar: req.QueryParams().Get("bar"),
			Baz: req.QueryParams().Get("baz"),
		}

		b, err := json.Marshal(m)
		assert.NoError(t, err)

		return &handler.Response{StatusCode: http.StatusOK, Body: string(b)}
	}

	resp := &ResponseWriter{}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, WithMiddlewares(
		func(r *Request) (*Request, error) {
			r.SetQueryByName("bar", "3")

			return r, nil
		},
		func(r *Request) (*Request, error) {
			r.SetQueryByName("baz", "4")

			return r, nil
		},
	)).Run()(resp, req)

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, `{"foo":"foo","bar":"3","baz":"4"}`, string(resp.Body))
}

func TestMiddlewaresError(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		return &handler.Response{StatusCode: http.StatusOK}
	}

	resp := &ResponseWriter{}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, WithMiddlewares(
		func(r *Request) (*Request, error) {
			return r, serviceerror.BadRequest("something bad has happened")
		},
	)).Run()(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, `{"error":{"id":"BAD_REQUEST","code":"BAD_REQUEST","message":"something bad has happened"}}`, string(resp.Body))
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

	resp := &ResponseWriter{}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, WithInterceptors(
		func(w *ResponseWriter) error {
			m := &metasyntactic{}
			err := json.Unmarshal([]byte(w.Body), m)
			assert.NoError(t, err)

			m.Bar = "4"

			b, err := json.Marshal(m)
			assert.NoError(t, err)

			w.Write(b)

			return nil
		},
	)).Run()(resp, req)

	assert.Equal(t, http.StatusOK, resp.Status)
	assert.Equal(t, `{"foo":"foo","bar":"4","baz":""}`, string(resp.Body))
}

func TestInterceptorsError(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		m := &metasyntactic{
			Foo: "foo",
		}

		b, err := json.Marshal(m)
		assert.NoError(t, err)

		return &handler.Response{StatusCode: http.StatusOK, Body: string(b)}
	}

	resp := &ResponseWriter{}
	req := &http.Request{
		Method: http.MethodGet,
		URL:    &url.URL{},
	}
	New(testHandler, WithInterceptors(
		func(w *ResponseWriter) error {
			return serviceerror.BadRequest("something bad has happened")
		},
	)).Run()(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Equal(t, `{"error":{"id":"BAD_REQUEST","code":"BAD_REQUEST","message":"something bad has happened"}}`, string(resp.Body))
}
