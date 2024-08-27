package handler

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/itsoneiota/lambda-handlers/pkg/test"
	"github.com/stretchr/testify/assert"
)

type reponseWriter struct {
	Headers http.Header
	Body    []byte
	Status  int
}

func (r *reponseWriter) Header() http.Header {
	if r.Headers == nil {
		r.Headers = http.Header{}
	}

	return r.Headers
}

func (r *reponseWriter) Write(body []byte) (int, error) {
	r.Body = body

	return len(body), nil
}

func (r *reponseWriter) WriteHeader(status int) {
	r.Status = status
}

type Model struct {
	Success bool `json:"success"`
}

func TestBuildResponder(t *testing.T) {
	body := "model"
	code := 200
	headers := http.Header{}
	headers.Set("default", "header")

	l := test.NewNullLogger()
	slog.SetDefault(l)

	hand := NewResponseHandler()
	res := &reponseWriter{
		Headers: headers,
	}

	err := hand.BuildResponderWithHeader(res, code, body, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, res.Status)
	assert.Equal(t, body, string(res.Body))
	assert.Equal(t, "header", res.Headers.Get("Default"))
}

func TestBuildResponseWithHeader_Empty(t *testing.T) {
	code := 200
	headers := http.Header{}
	headers.Set("default", "header")

	l := test.NewNullLogger()
	slog.SetDefault(l)

	hand := NewResponseHandler()
	res := &reponseWriter{
		Headers: headers,
	}

	err := hand.BuildResponseWithHeader(res, code, nil, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, res.Status)
	assert.Equal(t, "", string(res.Body))
	assert.Equal(t, "header", res.Headers.Get("Default"))
}

func TestBuildResponseWithHeader(t *testing.T) {
	model := Model{
		Success: true,
	}

	code := 200
	headers := http.Header{}
	headers.Set("default", "header")

	l := test.NewNullLogger()
	slog.SetDefault(l)

	hand := NewResponseHandler()
	res := &reponseWriter{
		Headers: headers,
	}

	err := hand.BuildResponseWithHeader(res, code, model, nil)
	assert.NoError(t, err)

	assert.Equal(t, 200, res.Status)
	assert.Equal(t, `{"success":true}`, string(res.Body))
	assert.Equal(t, "header", res.Headers.Get("Default"))
}

func TestBuildResponseWithHeader_Multiple(t *testing.T) {
	model := Model{
		Success: true,
	}

	code := 200
	headers := http.Header{
		"Server-Timing": []string{
			"cdn-cache; desc=HIT",
			"edge; dur=1",
			"ak_p; desc=\"467247_400071605_276706062_672_15674_1_0\";dur=1",
		},
	}

	hand := NewResponseHandler()
	res := &reponseWriter{}

	err := hand.BuildResponseWithHeader(res, code, model, headers)
	assert.NoError(t, err)

	assert.Equal(t, 200, res.Status)
	assert.Equal(t, `{"success":true}`, string(res.Body))
	assert.Equal(t, 3, len(res.Headers["Server-Timing"]))
}

func TestBuildResponseWithHeader_Cookie(t *testing.T) {
	model := Model{
		Success: true,
	}

	code := 200
	headers := http.Header{
		"Set-Cookie": []string{
			"loggedIn=True; path=/; secure",
			"QueueITAccepted-SDFrts345E-V3_nativeapptesta=EventId%3Dnativeapptesta%26QueueId%3Dd8dce36b-6cae-4e7b-b4af-c9bf27a9fb7f%26RedirectType%3Dqueue%26IssueTime%3D1685627225%26Hip%3D4c554f91c51760dd1d6a7fab0bbc41f7fe3023401c14a14db2174b1a677b747b%26Hash%3Deb133680cd4f6a163d6305ef0760a8b6b3abcef4e3b5907d6ef5f4d41c6d9dd4; expires=Fri, 02 Jun 2023 13:47:05 GMT; path=/",
			"loggedIn=True; path=/; secure; SameSite=None;Secure",
			"QueueITAccepted-SDFrts345E-V3_nativeapptesta=EventId%3Dnativeapptesta%26QueueId%3Dd8dce36b-6cae-4e7b-b4af-c9bf27a9fb7f%26RedirectType%3Dqueue%26IssueTime%3D1685627225%26Hip%3D4c554f91c51760dd1d6a7fab0bbc41f7fe3023401c14a14db2174b1a677b747b%26Hash%3Deb133680cd4f6a163d6305ef0760a8b6b3abcef4e3b5907d6ef5f4d41c6d9dd4; expires=Fri, 02 Jun 2023 13:47:05 GMT; path=/",
		},
	}

	l := test.NewNullLogger()
	slog.SetDefault(l)

	hand := NewResponseHandler()
	res := &reponseWriter{}

	err := hand.BuildResponseWithHeader(res, code, model, headers)

	assert.NoError(t, err)
	assert.Equal(t, headers, res.Headers)
}
