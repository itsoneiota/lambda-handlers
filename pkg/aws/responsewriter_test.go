package aws

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResponseWriterSuite struct {
	suite.Suite
	headers http.Header
}

func (s *ResponseWriterSuite) SetupTest() {
	s.headers = http.Header{
		"Content-Type": []string{
			"application/json; charset=utf-8",
		},
	}
}

func (s *ResponseWriterSuite) TestNewResponseWriter() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Headers)

	r.WriteHeader(http.StatusOK)
	s.Equal(http.StatusOK, r.StatusCode)

	r.Write([]byte("foo"))
	s.Equal("foo", r.Body)
}

func (s *ResponseWriterSuite) TestErrorResponse() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Headers)

	r.WriteHeader(http.StatusBadRequest)
	s.Equal(http.StatusBadRequest, r.StatusCode)

	r.Write([]byte("Oops"))
	s.Equal("{\"error\":{\"id\":\"BAD_REQUEST\",\"code\":\"BAD_REQUEST\",\"message\":\"Oops\"}}", r.Body)
}

func (s *ResponseWriterSuite) TestWrapperString() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Headers)

	r.WriteHeader(http.StatusBadRequest)
	s.Equal(http.StatusBadRequest, r.StatusCode)

	r.Write([]byte("\"Oops\"\n"))
	s.Equal("{\"error\":{\"id\":\"BAD_REQUEST\",\"code\":\"BAD_REQUEST\",\"message\":\"Oops\"}}", r.Body)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResponseWriterSuite(t *testing.T) {
	suite.Run(t, new(ResponseWriterSuite))
}
