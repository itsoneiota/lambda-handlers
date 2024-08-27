package aws

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponseWriter(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	r := NewResponseWriter(headers)

	assert.IsType(t, &ResponseWriter{}, r)
	assert.NotEmpty(t, r.Headers)

	r.Write([]byte("foo"))
	assert.NotEmpty(t, r.Body)

	r.WriteHeader(http.StatusOK)
	assert.Equal(t, http.StatusOK, r.StatusCode)
}
