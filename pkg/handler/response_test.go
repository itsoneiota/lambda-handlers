package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Model struct {
	Success bool `json:"success"`
}

func TestNewResponse(t *testing.T) {
	body := "model"
	code := http.StatusOK

	res := NewResponse(code, body)
	assert.IsType(t, (*Response)(nil), res)

	assert.Equal(t, code, res.StatusCode)
	assert.Equal(t, body, res.Body)
}

func TestNewResponseWithModel(t *testing.T) {
	body := Model{
		Success: true,
	}
	code := 200

	res := NewResponse(code, body)
	assert.IsType(t, (*Response)(nil), res)

	b, err := json.Marshal(body)
	assert.NoError(t, err)

	assert.JSONEq(t, string(b), res.Body)
}

func TestBuildResponseWithHeader(t *testing.T) {
	res := NewResponse(
		http.StatusOK,
		Model{
			Success: true,
		},
		WithHeaders(
			http.Header{
				"default": []string{"header"},
			},
		),
	)

	assert.IsType(t, (*Response)(nil), res)
	assert.NotEmpty(t, res.Headers)
}

func TestBuildResponseWithHeader_Multiple(t *testing.T) {
	headers := http.Header{
		"Server-Timing": []string{
			"cdn-cache; desc=HIT",
			"edge; dur=1",
			"ak_p; desc=\"467247_400071605_276706062_672_15674_1_0\";dur=1",
		},
	}

	res := NewResponse(
		http.StatusOK,
		Model{
			Success: true,
		},
		WithHeaders(headers),
	)

	assert.IsType(t, (*Response)(nil), res)
	assert.Equal(t, headers, res.Headers)
}

func TestBuildResponseWithHeader_Cookie(t *testing.T) {
	headers := http.Header{
		"Set-Cookie": []string{
			"loggedIn=True; path=/; secure",
			"QueueITAccepted-SDFrts345E-V3_nativeapptesta=EventId%3Dnativeapptesta%26QueueId%3Dd8dce36b-6cae-4e7b-b4af-c9bf27a9fb7f%26RedirectType%3Dqueue%26IssueTime%3D1685627225%26Hip%3D4c554f91c51760dd1d6a7fab0bbc41f7fe3023401c14a14db2174b1a677b747b%26Hash%3Deb133680cd4f6a163d6305ef0760a8b6b3abcef4e3b5907d6ef5f4d41c6d9dd4; expires=Fri, 02 Jun 2023 13:47:05 GMT; path=/",
			"loggedIn=True; path=/; secure; SameSite=None;Secure",
			"QueueITAccepted-SDFrts345E-V3_nativeapptesta=EventId%3Dnativeapptesta%26QueueId%3Dd8dce36b-6cae-4e7b-b4af-c9bf27a9fb7f%26RedirectType%3Dqueue%26IssueTime%3D1685627225%26Hip%3D4c554f91c51760dd1d6a7fab0bbc41f7fe3023401c14a14db2174b1a677b747b%26Hash%3Deb133680cd4f6a163d6305ef0760a8b6b3abcef4e3b5907d6ef5f4d41c6d9dd4; expires=Fri, 02 Jun 2023 13:47:05 GMT; path=/",
		},
	}

	res := NewResponse(
		http.StatusOK,
		Model{
			Success: true,
		},
		WithHeaders(headers),
	)
	assert.IsType(t, (*Response)(nil), res)

	assert.Equal(t, headers, res.Headers)
}
