package handler

import (
	"net/http"
	"testing"

	"github.com/itsoneiota/lambda-handlers/pkg/fakers"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testHandler := func(ctx Contexter, req Requester) *Response {
		return &Response{StatusCode: http.StatusOK, Body: "foo"}
	}

	req := fakers.Request{}
	ctx := fakers.Context{}
	resp := New(testHandler).Run()(ctx, req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "foo", resp.Body)
}

func TestRunResponseError(t *testing.T) {
	testHandler := func(ctx Contexter, req Requester) *Response {
		return &Response{StatusCode: http.StatusNotFound, Body: "foo not found"}
	}

	req := fakers.Request{}
	ctx := fakers.Context{}
	resp := New(testHandler).Run()(ctx, req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "foo not found", resp.Body)
}
