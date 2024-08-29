package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseHandlerSuite struct {
	suite.Suite
	status  int
	body    string
	headers http.Header
	resp    *reponseWriter
	handler *ResponseHandler
}

func (s *ResponseHandlerSuite) SetupTest() {
	s.status = http.StatusOK
	s.body = "model"
	s.headers = http.Header{}
	s.headers.Add("default", "header")
	s.resp = &reponseWriter{
		Headers: s.headers,
	}
	s.handler = NewResponseHandler()
}

func (s *ResponseHandlerSuite) TestBuildResponder() {
	err := s.handler.BuildResponder(s.resp, s.status, s.body)
	s.NoError(err)

	s.Equal(s.status, s.resp.Status)
	s.Equal(s.body, string(s.resp.Body))
	s.Equal("header", s.resp.Headers.Get("default"))
}

func (s *ResponseHandlerSuite) TestBuildResponseWithHeader_Empty() {
	err := s.handler.BuildResponseWithHeader(s.resp, s.status, nil, nil)
	s.NoError(err)

	s.Equal(s.status, s.resp.Status)
	s.Equal("", string(s.resp.Body))
	s.Equal("header", s.resp.Headers.Get("default"))
}

func (s *ResponseHandlerSuite) TestBuildResponseWithHeader() {
	model := Model{
		Success: true,
	}

	err := s.handler.BuildResponseWithHeader(s.resp, s.status, model, nil)
	s.NoError(err)

	s.Equal(s.status, s.resp.Status)
	s.Equal(`{"success":true}`, string(s.resp.Body))
	s.Equal("header", s.resp.Headers.Get("default"))
}

func (s *ResponseHandlerSuite) TestBuildResponseWithHeader_Multiple() {
	model := Model{
		Success: true,
	}

	headers := http.Header{
		"Server-Timing": []string{
			"cdn-cache; desc=HIT",
			"edge; dur=1",
			"ak_p; desc=\"467247_400071605_276706062_672_15674_1_0\";dur=1",
		},
	}

	err := s.handler.BuildResponseWithHeader(s.resp, s.status, model, headers)
	s.NoError(err)

	s.Equal(s.status, s.resp.Status)
	s.Equal(`{"success":true}`, string(s.resp.Body))
	s.Equal(3, len(s.resp.Headers["Server-Timing"]))
}

func (s *ResponseHandlerSuite) TestBuildResponseWithHeader_Cookie() {
	model := Model{
		Success: true,
	}

	headers := http.Header{
		"Set-Cookie": []string{
			"loggedIn=True; path=/; secure",
			"QueueITAccepted-SDFrts345E-V3_nativeapptesta=EventId%3Dnativeapptesta%26QueueId%3Dd8dce36b-6cae-4e7b-b4af-c9bf27a9fb7f%26RedirectType%3Dqueue%26IssueTime%3D1685627225%26Hip%3D4c554f91c51760dd1d6a7fab0bbc41f7fe3023401c14a14db2174b1a677b747b%26Hash%3Deb133680cd4f6a163d6305ef0760a8b6b3abcef4e3b5907d6ef5f4d41c6d9dd4; expires=Fri, 02 Jun 2023 13:47:05 GMT; path=/",
			"loggedIn=True; path=/; secure; SameSite=None;Secure",
			"QueueITAccepted-SDFrts345E-V3_nativeapptesta=EventId%3Dnativeapptesta%26QueueId%3Dd8dce36b-6cae-4e7b-b4af-c9bf27a9fb7f%26RedirectType%3Dqueue%26IssueTime%3D1685627225%26Hip%3D4c554f91c51760dd1d6a7fab0bbc41f7fe3023401c14a14db2174b1a677b747b%26Hash%3Deb133680cd4f6a163d6305ef0760a8b6b3abcef4e3b5907d6ef5f4d41c6d9dd4; expires=Fri, 02 Jun 2023 13:47:05 GMT; path=/",
		},
	}

	err := s.handler.BuildResponseWithHeader(s.resp, s.status, model, headers)
	s.NoError(err)

	headers.Add("default", "header")
	s.Equal(headers, s.resp.Headers)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResponseHandlerSuite(t *testing.T) {
	suite.Run(t, new(ResponseHandlerSuite))
}
