package aws

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/itsoneiota/lambda-handlers/pkg/handler"
	"github.com/stretchr/testify/assert"
)

func TestEncodeHeaders(t *testing.T) {
	h := http.Header{
		"Content-Type": []string{
			"application/json; charset=utf-8",
		},
		"Test": []string{
			"foo",
			"bar",
			"bar",
		},
	}

	expect := map[string]string{
		"Content-Type": "application/json; charset=utf-8",
		"Test":         "foo; bar",
	}

	assert.Equal(t, expect, encodeHeaders(h))
}

func TestHandle(t *testing.T) {
	type metasyntactic struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		m := &metasyntactic{Foo: "handler"}
		b, err := json.Marshal(m)
		assert.NoError(t, err)

		return &handler.Response{StatusCode: http.StatusOK, Body: string(b)}
	}

	h := &Handler{
		handler: handler.New(testHandler),
		Opt: &Opt{
			interceptors: []handler.Interceptor{
				func(
					_ handler.Contexter,
					_ handler.Requester,
					r *handler.Response,
				) *handler.Response {
					m := &metasyntactic{}
					err := json.Unmarshal([]byte(r.Body), m)
					assert.NoError(t, err)

					m.Bar = "interceptor 1"

					b, err := json.Marshal(m)
					assert.NoError(t, err)

					r.Body = string(b)

					return r
				},
				func(
					_ handler.Contexter,
					_ handler.Requester,
					r *handler.Response,
				) *handler.Response {
					m := &metasyntactic{}
					err := json.Unmarshal([]byte(r.Body), m)
					assert.NoError(t, err)

					m.Bar = "interceptor 2"

					b, err := json.Marshal(m)
					assert.NoError(t, err)

					r.Body = string(b)

					return r
				},
			},
		},
	}

	req := &events.APIGatewayProxyRequest{}
	resp, err := handle(h)(req)
	assert.NoError(t, err)

	assert.JSONEq(t, `{"foo":"handler","bar":"interceptor 2"}`, resp.Body)

}

func TestError(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		return &handler.Response{StatusCode: http.StatusBadRequest, Body: "Oh no!"}
	}

	h := &Handler{
		handler: handler.New(testHandler),
		Opt: &Opt{
			interceptors: []handler.Interceptor{
				func(
					_ handler.Contexter,
					_ handler.Requester,
					r *handler.Response,
				) *handler.Response {
					r.Body = "Oh yes!"

					return r
				},
			},
		},
	}

	req := &events.APIGatewayProxyRequest{}
	resp, err := handle(h)(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Oh no!", resp.Body)
}

func TestPanic(t *testing.T) {
	testHandler := func(ctx handler.Contexter, req handler.Requester) *handler.Response {
		panic("test")
	}

	h := &Handler{
		handler: handler.New(testHandler),
		Opt:     &Opt{},
	}

	req := &events.APIGatewayProxyRequest{}
	resp, err := handle(h)(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.JSONEq(t, `{"error": {"code":"INTERNAL_SERVER_ERROR", "id":"INTERNAL_SERVER_ERROR", "message":"Internal Server Error"}}`, resp.Body)

}
