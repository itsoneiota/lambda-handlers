package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	testHandler := func(ctx Contexter, req Requester) *Response {
		body := map[string]any{
			"foo": ctx.Value("foo"),
			"bar": ctx.Value("bar"),
		}
		b, err := json.Marshal(body)
		assert.NoError(t, err)

		return &Response{StatusCode: http.StatusOK, Body: string(b)}
	}

	testMiddlewareOne := func(next HandlerFunc) HandlerFunc {
		return func(ctx Contexter, req Requester) *Response {
			ctx.SetValue("foo", 1)
			return next(ctx, req)
		}
	}

	testMiddlewareTwo := func(next HandlerFunc) HandlerFunc {
		return func(ctx Contexter, req Requester) *Response {
			ctx.SetValue("bar", 2)
			return next(ctx, req)
		}
	}

	req := fakeRequest{}
	ctx := fakeContext{
		values: map[string]any{},
	}
	resp := New(testHandler, WithHeaders(http.Header{"Accept": []string{"application/json"}})).
		Middlewares(testMiddlewareOne, testMiddlewareTwo).
		Run()(ctx, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"bar":2,"foo":1}`, resp.Body)
	assert.NotEmpty(t, resp.Headers)
}
