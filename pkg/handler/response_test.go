package handler

import (
	"net/http"
	"testing"

	"github.com/itsoneiota/lambda-handlers/pkg/handler/mocks"
	"github.com/stretchr/testify/assert"
)

type Model struct {
	Success bool `json:"success"`
}

func TestBuildResponder(t *testing.T) {
	body := "model"
	code := 200
	headers := http.Header{}
	headers.Set("default", "header")

	l := mocks.NewLogger(t)

	hand := NewResponseHandler(l, headers)

	res, err := hand.BuildResponderWithHeader(code, body, nil)

	assert.NoError(t, err)
	assert.IsType(t, (*Response)(nil), res)
}

func TestBuildResponseWithHeader_Empty(t *testing.T) {
	code := 200
	headers := http.Header{}
	headers.Set("default", "header")

	l := mocks.NewLogger(t)

	hand := NewResponseHandler(l, headers)

	res, err := hand.BuildResponseWithHeader(code, nil, nil)

	assert.NoError(t, err)
	assert.IsType(t, (*Response)(nil), res)
}

func TestBuildResponseWithHeader(t *testing.T) {
	model := Model{
		Success: true,
	}

	code := 200
	headers := http.Header{}
	headers.Set("default", "header")

	l := mocks.NewLogger(t)

	hand := NewResponseHandler(l, headers)

	res, err := hand.BuildResponseWithHeader(code, model, nil)

	assert.NoError(t, err)
	assert.IsType(t, (*Response)(nil), res)
}
