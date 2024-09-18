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
	s.NotEmpty(r.Header())

	r.WriteHeader(http.StatusOK)
	s.Equal(http.StatusOK, r.StatusCode)

	r.Write([]byte("foo"))
	s.Equal("foo", r.Body)
}

func (s *ResponseWriterSuite) TestErrorResponse() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Header())

	r.WriteHeader(http.StatusBadRequest)
	s.Equal(http.StatusBadRequest, r.StatusCode)

	r.Write([]byte("Oops"))
	s.Equal("{\"error\":{\"id\":\"BAD_REQUEST\",\"code\":\"BAD_REQUEST\",\"message\":\"Oops\"}}", r.Body)
}

func (s *ResponseWriterSuite) TestWrapperString() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Header())

	r.WriteHeader(http.StatusBadRequest)
	s.Equal(http.StatusBadRequest, r.StatusCode)

	r.Write([]byte("\"Oops\"\n"))
	s.Equal("{\"error\":{\"id\":\"BAD_REQUEST\",\"code\":\"BAD_REQUEST\",\"message\":\"Oops\"}}", r.Body)
}

func (s *ResponseWriterSuite) TestAddHeader() {
	r := NewResponseWriter(s.headers)

	s.IsType(&ResponseWriter{}, r)
	s.NotEmpty(r.Header())

	r.Header().Add("foo", "bar")
	s.Equal(2, len(r.Header()))
	s.Equal("bar", r.Header().Get("foo"))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestResponseWriterSuite(t *testing.T) {
	suite.Run(t, new(ResponseWriterSuite))
}
