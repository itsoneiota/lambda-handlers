package aws

import (
	"net/http"
	"testing"

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
