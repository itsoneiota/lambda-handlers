package mux

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContexterInterface(t *testing.T) {
	ip := "SourceIP"
	r := &http.Request{
		RemoteAddr: "SourceIP",
		Method:     http.MethodGet,
	}

	ctx := Context{Request: r}

	assert.Equal(t, ip, ctx.SourceIP())
	assert.IsType(t, int64(1), ctx.UnixNow())
	assert.IsType(t, http.MethodGet, ctx.HttpMethod())
	assert.IsType(t, "dev", ctx.Stage())
}